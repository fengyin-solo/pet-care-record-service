package service

import (
	"errors"

	"petsmanagement/internal/model"
)

type ReminderSender interface{ Send(string) error }

func (s *Service) DeliverReminder(key string, sender ReminderSender) error {
	delivery := model.ReminderDelivery{Key: key, Status: "sending"}
	for delivery.Attempts < 2 {
		delivery.Attempts++
		err := sender.Send(key)
		if err == nil {
			delivery.Status = "delivered"
			s.store.PutReminderDelivery(delivery)
			return nil
		}
		normalized := model.NormalizeDeliveryError(err)
		delivery.Status = "failed"
		s.store.PutReminderDelivery(delivery)
		if !model.IsTemporaryDelivery(normalized) {
			return normalized
		}
	}
	return errors.New("reminder delivery exhausted")
}

func (s *Service) ReminderDelivery(key string) (model.ReminderDelivery, bool) {
	return s.store.ReminderDelivery(key)
}
