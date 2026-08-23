package service

import (
	"errors"

	"petsmanagement/internal/adapter"
	"petsmanagement/internal/model"
)

func (s *Service) UpdateHealthWithRetry(petID, status string, publisher *adapter.HealthPublisher) error {
	var firstErr error
	for attempt := 1; attempt <= 2; attempt++ {
		update := model.NewHealthUpdate(petID, status, attempt)
		s.store.ApplyHealthUpdate(update)
		if err := publisher.Publish(update.Key); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			if errors.Is(err, adapter.ErrPublisherBusy) {
				continue
			}
			return err
		}
		return firstErr
	}
	return firstErr
}

func (s *Service) HealthUpdateHistory() []model.HealthUpdate { return s.store.HealthUpdateHistory() }
