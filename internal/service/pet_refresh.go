package service

import (
	"context"
	"sync"

	"petsmanagement/internal/refreshschedule"
	"petsmanagement/internal/refreshworker"
)

type RefreshCoordinator struct {
	worker *refreshworker.Worker
	wg     sync.WaitGroup
}

func NewRefreshCoordinator(client refreshworker.Client) *RefreshCoordinator {
	return &RefreshCoordinator{worker: refreshworker.New(client)}
}

func (c *RefreshCoordinator) Refresh(ctx context.Context, petID string) error {
	done := make(chan error, 1)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		done <- refreshschedule.Run(context.Background(), petID, c.worker)
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *RefreshCoordinator) Shutdown() { c.wg.Wait() }
