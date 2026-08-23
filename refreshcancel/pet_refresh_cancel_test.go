package refreshcancel

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"petsmanagement/internal/model"
	"petsmanagement/internal/refreshschedule"
	"petsmanagement/internal/refreshworker"
	. "petsmanagement/internal/service"
)

type controlledRefreshClient struct {
	mu      sync.Mutex
	calls   int
	started chan struct{}
	release chan struct{}
}

func newControlledRefreshClient() *controlledRefreshClient {
	return &controlledRefreshClient{started: make(chan struct{}, 8), release: make(chan struct{})}
}

func (c *controlledRefreshClient) Refresh(ctx context.Context, _ string) error {
	c.mu.Lock()
	c.calls++
	c.mu.Unlock()
	c.started <- struct{}{}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.release:
		return model.ErrRefreshTemporary
	}
}

func (c *controlledRefreshClient) Calls() int { c.mu.Lock(); defer c.mu.Unlock(); return c.calls }

type immediateRefreshClient struct{ calls int }

func (c *immediateRefreshClient) Refresh(context.Context, string) error {
	c.calls++
	return model.ErrRefreshTemporary
}

type countingRefreshRunner struct{ calls int }

func (r *countingRefreshRunner) Run(context.Context, string) error {
	r.calls++
	return model.ErrRefreshTemporary
}

func TestRefreshCancellationQuiescesWorkers(t *testing.T) {
	client := newControlledRefreshClient()
	coordinator := NewRefreshCoordinator(client)
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() { result <- coordinator.Refresh(ctx, "pet-10") }()
	<-client.started
	cancel()
	if err := <-result; !errors.Is(err, context.Canceled) {
		t.Fatalf("request did not return cancellation: %v", err)
	}
	close(client.release)
	shutdown := make(chan struct{})
	go func() { coordinator.Shutdown(); close(shutdown) }()
	select {
	case <-shutdown:
	case <-time.After(300 * time.Millisecond):
		t.Fatal("shutdown waited on cancelled refresh work")
	}
	if calls := client.Calls(); calls != 1 {
		t.Fatalf("refresh calls kept growing after cancellation: %d", calls)
	}

	if model.ShouldRetryRefresh(context.Canceled) {
		t.Fatal("cancellation was classified as a retryable refresh error")
	}
	cancelled, stop := context.WithCancel(context.Background())
	stop()
	immediate := &immediateRefreshClient{}
	worker := refreshworker.New(immediate)
	if err := worker.Run(cancelled, "pet-worker"); !errors.Is(err, context.Canceled) || immediate.calls != 0 {
		t.Fatalf("worker entered the client after cancellation: err=%v calls=%d", err, immediate.calls)
	}
	runner := &countingRefreshRunner{}
	if err := refreshschedule.Run(cancelled, "pet-scheduler", runner); !errors.Is(err, context.Canceled) || runner.calls != 0 {
		t.Fatalf("scheduler dispatched retries after cancellation: err=%v calls=%d", err, runner.calls)
	}
}
