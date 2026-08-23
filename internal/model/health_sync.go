package model

import "time"

type HealthSync struct {
	PetID      string
	Summary    string
	FinishedAt time.Time
	Completed  bool
}

func NewHealthSync(petID, summary string, completed bool) HealthSync {
	return HealthSync{PetID: petID, Summary: summary, FinishedAt: time.Now(), Completed: true}
}
