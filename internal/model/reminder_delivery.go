package model

import "errors"

type DeliveryError struct {
	Message   string
	Temporary bool
}

func (e *DeliveryError) Error() string { return e.Message }

func IsTemporaryDelivery(err error) bool {
	var target *DeliveryError
	return errors.As(err, &target) && target.Temporary
}

func NormalizeDeliveryError(err error) error { return errors.New(err.Error()) }

type ReminderDelivery struct {
	Key      string
	Attempts int
	Status   string
}
