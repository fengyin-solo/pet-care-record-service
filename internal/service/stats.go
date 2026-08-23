package service

import (
	"time"

	"petsmanagement/internal/model"
)

// StatsOverview 宠物管理概览统计。
type StatsOverview struct {
	TotalOwners         int            `json:"total_owners"`
	TotalPets           int            `json:"total_pets"`
	PetsBySpecies       map[string]int `json:"pets_by_species"`
	PetsByHealthStatus  map[string]int `json:"pets_by_health_status"`
	PetsArchived        int            `json:"pets_archived"`
	TotalMedicalRecords int            `json:"total_medical_records"`
	TotalVaccines       int            `json:"total_vaccines"`
	PendingFollowUps    int            `json:"pending_follow_ups"`
	DueVaccinesIn30Days int            `json:"due_vaccines_in_30_days"`
}

// Overview 计算整体概览统计。
func (s *Service) Overview() *StatsOverview {
	o := &StatsOverview{
		PetsBySpecies:      make(map[string]int),
		PetsByHealthStatus: make(map[string]int),
	}

	o.TotalOwners = len(s.store.ListOwners())
	o.TotalMedicalRecords = len(s.store.ListMedicalRecords())
	o.TotalVaccines = len(s.store.ListVaccines())

	for _, p := range s.store.ListPets() {
		o.TotalPets++
		o.PetsBySpecies[p.Species]++
		o.PetsByHealthStatus[p.HealthStatus]++
		if p.Status == model.PetStatusArchived {
			o.PetsArchived++
		}
	}

	for _, f := range s.store.ListFollowUps() {
		if f.Status == model.FollowUpPending {
			o.PendingFollowUps++
		}
	}

	dueBefore := time.Now().AddDate(0, 0, 30)
	for _, v := range s.store.ListVaccines() {
		if !v.NextDueAt.IsZero() && v.NextDueAt.Before(dueBefore) {
			o.DueVaccinesIn30Days++
		}
	}

	return o
}
