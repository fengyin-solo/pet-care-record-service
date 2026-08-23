package model

import (
	"strings"
	"time"
)

// MedicalRecord 宠物就诊记录。
type MedicalRecord struct {
	ID             string    `json:"id"`
	PetID          string    `json:"pet_id"`
	VetName        string    `json:"vet_name"`
	Clinic         string    `json:"clinic"`
	Diagnosis      string    `json:"diagnosis"`
	Treatment      string    `json:"treatment"`
	VisitAt        time.Time `json:"visit_at"`
	FollowUpNeeded bool      `json:"follow_up_needed"`
	CreatedAt      time.Time `json:"created_at"`
}

// Validate 校验就诊记录字段并规范化。
func (m *MedicalRecord) Validate() error {
	m.PetID = strings.TrimSpace(m.PetID)
	m.VetName = strings.TrimSpace(m.VetName)
	m.Clinic = strings.TrimSpace(m.Clinic)
	m.Diagnosis = strings.TrimSpace(m.Diagnosis)
	m.Treatment = strings.TrimSpace(m.Treatment)
	if m.PetID == "" {
		return NewValidationError("pet_id", "宠物 ID 不能为空")
	}
	if m.VetName == "" {
		return NewValidationError("vet_name", "兽医姓名不能为空")
	}
	if m.VisitAt.IsZero() {
		m.VisitAt = time.Now()
	}
	return nil
}

// MedicalRecordFilter 就诊记录查询筛选条件。
type MedicalRecordFilter struct {
	PetID  string
	Clinic string
}

// Match 判断就诊记录是否命中筛选条件。
func (f MedicalRecordFilter) Match(m *MedicalRecord) bool {
	if f.PetID != "" && m.PetID != f.PetID {
		return false
	}
	if f.Clinic != "" && m.Clinic != f.Clinic {
		return false
	}
	return true
}
