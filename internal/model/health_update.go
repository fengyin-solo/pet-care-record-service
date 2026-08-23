package model

import "fmt"

type HealthUpdate struct {
	PetID  string
	Status string
	Key    string
}

func NewHealthUpdate(petID, status string, attempt int) HealthUpdate {
	return HealthUpdate{PetID: petID, Status: status, Key: fmt.Sprintf("%s-%d", petID, attempt)}
}
