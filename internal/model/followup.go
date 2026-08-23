package model

import (
	"strings"
	"time"
)

// 回访方式与状态常量。
const (
	FollowUpPhone = "phone"
	FollowUpVisit = "visit"

	FollowUpPending   = "pending"
	FollowUpDone      = "done"
	FollowUpCancelled = "cancelled"
)

// FollowUp 宠物回访记录。
type FollowUp struct {
	ID          string     `json:"id"`
	PetID       string     `json:"pet_id"`
	OwnerID     string     `json:"owner_id"`
	Method      string     `json:"method"`
	Status      string     `json:"status"`
	ScheduledAt time.Time  `json:"scheduled_at"`
	Note        string     `json:"note"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Validate 校验回访记录字段并规范化。
func (f *FollowUp) Validate() error {
	f.PetID = strings.TrimSpace(f.PetID)
	f.OwnerID = strings.TrimSpace(f.OwnerID)
	f.Method = strings.TrimSpace(f.Method)
	f.Note = strings.TrimSpace(f.Note)
	if f.PetID == "" {
		return NewValidationError("pet_id", "宠物 ID 不能为空")
	}
	if f.Method == "" {
		f.Method = FollowUpPhone
	}
	if f.Method != FollowUpPhone && f.Method != FollowUpVisit {
		return NewValidationError("method", "回访方式不合法")
	}
	if f.Status == "" {
		f.Status = FollowUpPending
	}
	if f.Status != FollowUpPending && f.Status != FollowUpDone && f.Status != FollowUpCancelled {
		return NewValidationError("status", "回访状态不合法")
	}
	if f.ScheduledAt.IsZero() {
		f.ScheduledAt = time.Now()
	}
	return nil
}

// followUpTransitions 回访状态机：pending→done、pending→cancelled。
var followUpTransitions = map[string]map[string]bool{
	FollowUpPending: {FollowUpDone: true, FollowUpCancelled: true},
}

// CanTransitionFollowUpStatus 判断回访状态是否允许从 from 流转到 to。
func CanTransitionFollowUpStatus(from, to string) bool {
	if m, ok := followUpTransitions[from]; ok {
		return m[to]
	}
	return false
}

// FollowUpFilter 回访记录查询筛选条件。
type FollowUpFilter struct {
	PetID  string
	Status string
}

// Match 判断回访记录是否命中筛选条件。
func (f FollowUpFilter) Match(fu *FollowUp) bool {
	if f.PetID != "" && fu.PetID != f.PetID {
		return false
	}
	if f.Status != "" && fu.Status != f.Status {
		return false
	}
	return true
}
