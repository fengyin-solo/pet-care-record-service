package store

import "petsmanagement/internal/model"

func (s *MemoryStore) SavePetDraft(draft model.PetDraft) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.petDrafts[draft.Key] = draft
	return true
}

func (s *MemoryStore) PetDraft(key string) (model.PetDraft, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	draft, ok := s.petDrafts[key]
	return draft, ok
}
