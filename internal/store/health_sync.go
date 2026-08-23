package store

import "petsmanagement/internal/model"

func (s *MemoryStore) AddHealthSync(sync model.HealthSync) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.healthSyncs = append(s.healthSyncs, sync)
	return true
}

func (s *MemoryStore) ListHealthSyncs() []model.HealthSync {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.HealthSync(nil), s.healthSyncs...)
}
