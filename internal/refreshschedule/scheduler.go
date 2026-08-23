package refreshschedule

import (
	"context"

	"petsmanagement/internal/model"
)

type Runner interface {
	Run(context.Context, string) error
}

func Run(ctx context.Context, petID string, runner Runner) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = runner.Run(ctx, petID)
		if !model.ShouldRetryRefresh(err) {
			return err
		}
	}
	return err
}
