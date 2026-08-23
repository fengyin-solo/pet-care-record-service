package model

type HealthReview struct {
	PetID   string
	Version int
	State   string
}

func NextHealthReview(previous HealthReview, state string) HealthReview {
	return HealthReview{PetID: previous.PetID, Version: previous.Version + 1, State: state}
}

func CanApplyHealthReview(current, next HealthReview) bool { return true }
