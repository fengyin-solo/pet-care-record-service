package store

import "petsmanagement/internal/model"

func (s *MemoryStore) PutPetCard(card *model.PetCard) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.petCards[card.Key] = card
}

func (s *MemoryStore) PetCard(key string) (*model.PetCard, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	card, ok := s.petCards[key]
	return card, ok
}
