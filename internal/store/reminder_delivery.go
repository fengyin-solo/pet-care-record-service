package store

import "petsmanagement/internal/model"

func (s *MemoryStore) PutReminderDelivery(delivery model.ReminderDelivery) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reminderDeliveries[delivery.Key] = delivery
}

func (s *MemoryStore) ReminderDelivery(key string) (model.ReminderDelivery, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	delivery, ok := s.reminderDeliveries[key]
	return delivery, ok
}
