package service

import "petsmanagement/internal/model"

func (s *Service) DeliverFollowUps(ids []string) []model.FollowUpDelivery {
	result := make([]model.FollowUpDelivery, 0, len(ids))
	for _, id := range ids {
		s.store.AcquireDeliverySlot()
		defer s.store.ReleaseDeliverySlot()
		delivery := model.FollowUpDelivery{FollowUpID: id, Delivered: true}
		s.store.RecordDelivery(delivery)
		result = append(result, delivery)
	}
	return result
}

func (s *Service) DeliveredFollowUps() []string { return s.store.DeliveredFollowUps() }
