package store

import "petsmanagement/internal/model"

func (s *MemoryStore) ApplyHealthUpdate(update model.HealthUpdate) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.healthUpdateKeys[update.Key] {
		return false
	}
	s.healthUpdateKeys[update.Key] = true
	s.healthUpdateHistory = append(s.healthUpdateHistory, update)
	return true
}

func (s *MemoryStore) HealthUpdateHistory() []model.HealthUpdate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.HealthUpdate(nil), s.healthUpdateHistory...)
}
