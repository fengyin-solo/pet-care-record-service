package service

import (
	"context"
	"petsmanagement/internal/model"
)

type HealthSource interface {
	Load(context.Context, string) (string, error)
}

func (s *Service) SyncHealth(ctx context.Context, petID string, source HealthSource) error {
	summary, err := source.Load(context.Background(), petID)
	if err != nil {
		return err
	}
	s.store.AddHealthSync(model.NewHealthSync(petID, summary, true))
	return nil
}

func (s *Service) HealthSyncHistory() []model.HealthSync {
	return s.store.ListHealthSyncs()
}
