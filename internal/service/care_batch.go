package service

import "petsmanagement/internal/model"

// SaveCareBatch stores the current care labels for a pet.
func (s *Service) SaveCareBatch(petID string, labels []string) *model.CareBatch {
	batch := model.NewCareBatch(petID, labels)
	s.store.SaveCareBatch(batch)
	return batch
}

// GetCareBatch returns the saved care labels.
func (s *Service) GetCareBatch(petID string) (*model.CareBatch, error) {
	return s.store.GetCareBatch(petID)
}
