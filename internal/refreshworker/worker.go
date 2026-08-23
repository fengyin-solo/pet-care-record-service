package refreshworker

import (
	"context"

	"petsmanagement/internal/model"
)

type Client interface {
	Refresh(context.Context, string) error
}

type Worker struct{ client Client }

func New(client Client) *Worker { return &Worker{client: client} }

// Run 执行单次下游刷新调用。
// 进入 worker 前若上下文已取消，立即返回，保证下游调用为 0。
// 将 ctx 透传给下游，使取消信号能及时传播。
func (w *Worker) Run(ctx context.Context, petID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := w.client.Refresh(ctx, petID)
	if err != nil && model.ShouldRetryRefresh(err) {
		return model.ErrRefreshTemporary
	}
	return err
}
