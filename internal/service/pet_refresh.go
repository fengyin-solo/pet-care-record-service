// Package service 实现业务逻辑层。
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

// Refresh 异步执行宠物资料刷新。
// 将请求 ctx 透传给后台 goroutine，使取消信号沿
// Coordinator → Scheduler → Worker → Client 整条链传播。
// ctx 取消后 Refresh 立即返回 ctx.Err()，
// 后台 goroutine 也会因 ctx 取消而及时退出，不阻塞 Shutdown。
func (c *RefreshCoordinator) Refresh(ctx context.Context, petID string) error {
	done := make(chan error, 1)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		done <- refreshschedule.Run(ctx, petID, c.worker)
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *RefreshCoordinator) Shutdown() { c.wg.Wait() }
