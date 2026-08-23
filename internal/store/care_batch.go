package store

import "petsmanagement/internal/model"

func (s *MemoryStore) SaveCareBatch(batch *model.CareBatch) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.careBatches[batch.PetID] = batch
}

func (s *MemoryStore) GetCareBatch(petID string) (*model.CareBatch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	batch, ok := s.careBatches[petID]
	if !ok {
		return nil, ErrNotFound
	}
	return batch, nil
}
