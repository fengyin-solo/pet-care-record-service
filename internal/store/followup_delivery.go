package store

import "petsmanagement/internal/model"

func (s *MemoryStore) AcquireDeliverySlot() { s.deliverySlots <- struct{}{} }

func (s *MemoryStore) ReleaseDeliverySlot() { <-s.deliverySlots }

func (s *MemoryStore) RecordDelivery(delivery model.FollowUpDelivery) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deliveredFollowUps = append(s.deliveredFollowUps, delivery.FollowUpID)
}

func (s *MemoryStore) DeliveredFollowUps() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.deliveredFollowUps...)
}
