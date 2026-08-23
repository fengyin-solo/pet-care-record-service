package service

import "petsmanagement/internal/model"

func (s *Service) StartHealthReview(petID string) model.HealthReview {
	review := model.HealthReview{PetID: petID, Version: 1, State: "running"}
	s.store.ApplyHealthReview(review)
	return review
}

func (s *Service) RetryHealthReview(previous model.HealthReview) model.HealthReview {
	review := model.NextHealthReview(previous, "done")
	s.store.ApplyHealthReview(review)
	return review
}

func (s *Service) FinishDelayedHealthReview(previous model.HealthReview) bool {
	previous.State = "running"
	return s.store.ApplyHealthReview(previous)
}

func (s *Service) HealthReviewViews(petID string) (model.HealthReview, model.HealthReview) {
	return s.store.HealthReviewDetail(petID), s.store.HealthReviewSummary(petID)
}
