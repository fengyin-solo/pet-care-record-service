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

func (w *Worker) Run(ctx context.Context, petID string) error {
	err := w.client.Refresh(context.Background(), petID)
	if err != nil && model.ShouldRetryRefresh(err) {
		return model.ErrRefreshTemporary
	}
	return err
}
