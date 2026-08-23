package model

import (
	"strings"
	"time"
)

// 宠物性别与健康状态、在档状态常量。
const (
	PetGenderMale    = "male"
	PetGenderFemale  = "female"
	PetGenderUnknown = "unknown"

	HealthHealthy    = "healthy"
	HealthSick       = "sick"
	HealthRecovering = "recovering"
	HealthTreatment  = "treatment"

	PetStatusActive   = "active"
	PetStatusArchived = "archived"
)

// Pet 宠物档案。
type Pet struct {
	ID           string    `json:"id"`
	OwnerID      string    `json:"owner_id"`
	Name         string    `json:"name"`
	Species      string    `json:"species"`
	Breed        string    `json:"breed"`
	Gender       string    `json:"gender"`
	BirthDate    time.Time `json:"birth_date"`
	HealthStatus string    `json:"health_status"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Validate 校验宠物字段并规范化。
func (p *Pet) Validate() error {
	p.OwnerID = strings.TrimSpace(p.OwnerID)
	p.Name = strings.TrimSpace(p.Name)
	p.Species = strings.TrimSpace(p.Species)
	p.Breed = strings.TrimSpace(p.Breed)
	p.Gender = strings.TrimSpace(p.Gender)
	p.HealthStatus = strings.TrimSpace(p.HealthStatus)

	if p.OwnerID == "" {
		return NewValidationError("owner_id", "主人 ID 不能为空")
	}
	if p.Name == "" {
		return NewValidationError("name", "宠物名字不能为空")
	}
	if p.Species == "" {
		return NewValidationError("species", "物种不能为空")
	}
	if p.Gender == "" {
		p.Gender = PetGenderUnknown
	}
	if p.Gender != PetGenderMale && p.Gender != PetGenderFemale && p.Gender != PetGenderUnknown {
		return NewValidationError("gender", "性别不合法")
	}
	if p.HealthStatus == "" {
		p.HealthStatus = HealthHealthy
	}
	if p.HealthStatus != HealthHealthy && p.HealthStatus != HealthSick &&
		p.HealthStatus != HealthRecovering && p.HealthStatus != HealthTreatment {
		return NewValidationError("health_status", "健康状态不合法")
	}
	if p.Status == "" {
		p.Status = PetStatusActive
	}
	if p.Status != PetStatusActive && p.Status != PetStatusArchived {
		return NewValidationError("status", "在档状态不合法")
	}
	return nil
}

// petTransitions 宠物状态机：仅允许 active→archived，归档不可逆。
var petTransitions = map[string]map[string]bool{
	PetStatusActive: {PetStatusArchived: true},
}

// CanTransitionPetStatus 判断宠物状态是否允许从 from 流转到 to。
func CanTransitionPetStatus(from, to string) bool {
	if m, ok := petTransitions[from]; ok {
		return m[to]
	}
	return false
}

// PetFilter 宠物查询筛选条件。
type PetFilter struct {
	OwnerID      string
	Species      string
	HealthStatus string
	Status       string
	Keyword      string
}

// Match 判断宠物是否命中筛选条件。
func (f PetFilter) Match(p *Pet) bool {
	if f.OwnerID != "" && p.OwnerID != f.OwnerID {
		return false
	}
	if f.Species != "" && p.Species != f.Species {
		return false
	}
	if f.HealthStatus != "" && p.HealthStatus != f.HealthStatus {
		return false
	}
	if f.Status != "" && p.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" &&
			!strings.Contains(strings.ToLower(p.Name), k) &&
			!strings.Contains(strings.ToLower(p.Breed), k) {
			return false
		}
	}
	return true
}
