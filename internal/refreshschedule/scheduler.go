package refreshschedule

import (
	"context"

	"petsmanagement/internal/model"
)

type Runner interface {
	Run(context.Context, string) error
}

// Run 执行刷新任务，最多重试 3 次。
// 进入前若上下文已取消，直接返回，不触发任何 runner 调用。
// 每次尝试后若上下文已取消/超时，立即退出，不再重试。
func Run(ctx context.Context, petID string, runner Runner) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = runner.Run(ctx, petID)
		if !model.ShouldRetryRefresh(err) {
			return err
		}
		// 上下文已取消或超时，不再重试
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return err
}
