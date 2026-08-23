package model

import (
	"strings"
	"time"
)

// Vaccine 宠物疫苗接种记录。
type Vaccine struct {
	ID             string    `json:"id"`
	PetID          string    `json:"pet_id"`
	VaccineName    string    `json:"vaccine_name"`
	BatchNo        string    `json:"batch_no"`
	AdministeredAt time.Time `json:"administered_at"`
	NextDueAt      time.Time `json:"next_due_at"`
	VetName        string    `json:"vet_name"`
	CreatedAt      time.Time `json:"created_at"`
}

// Validate 校验疫苗记录字段并规范化。
func (v *Vaccine) Validate() error {
	v.PetID = strings.TrimSpace(v.PetID)
	v.VaccineName = strings.TrimSpace(v.VaccineName)
	v.BatchNo = strings.TrimSpace(v.BatchNo)
	v.VetName = strings.TrimSpace(v.VetName)
	if v.PetID == "" {
		return NewValidationError("pet_id", "宠物 ID 不能为空")
	}
	if v.VaccineName == "" {
		return NewValidationError("vaccine_name", "疫苗名称不能为空")
	}
	if v.AdministeredAt.IsZero() {
		v.AdministeredAt = time.Now()
	}
	if !v.NextDueAt.IsZero() && v.NextDueAt.Before(v.AdministeredAt) {
		return NewValidationError("next_due_at", "下次到期时间不能早于接种时间")
	}
	return nil
}

// VaccineFilter 疫苗记录查询筛选条件。
type VaccineFilter struct {
	PetID     string
	DueBefore *time.Time
}

// Match 判断疫苗记录是否命中筛选条件。
func (f VaccineFilter) Match(v *Vaccine) bool {
	if f.PetID != "" && v.PetID != f.PetID {
		return false
	}
	if f.DueBefore != nil {
		if v.NextDueAt.IsZero() || v.NextDueAt.After(*f.DueBefore) {
			return false
		}
	}
	return true
}
