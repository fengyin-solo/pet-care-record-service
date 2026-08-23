package store

import "petsmanagement/internal/model"

func (s *MemoryStore) ApplyHealthReview(review model.HealthReview) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	currentDetail := s.healthReviewDetails[review.PetID]
	if !model.CanApplyHealthReview(currentDetail, review) {
		return false
	}
	s.healthReviewDetails[review.PetID] = review
	current := s.healthReviewSummaries[review.PetID]
	if review.Version > current.Version {
		s.healthReviewSummaries[review.PetID] = review
	}
	return true
}

func (s *MemoryStore) HealthReviewDetail(petID string) model.HealthReview {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.healthReviewDetails[petID]
}

func (s *MemoryStore) HealthReviewSummary(petID string) model.HealthReview {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.healthReviewSummaries[petID]
}
