package model

import "fmt"

// HealthUpdate 表示一次宠物健康状态变更事件。
type HealthUpdate struct {
	PetID  string
	Status string
	// Key 是幂等标识，仅由宠物 ID 与目标状态决定，不含尝试序号。
	// 同一宠物+状态的多次尝试共享同一个 Key，确保下游可做幂等去重。
	Key string
}

// NewHealthUpdate 构造健康状态变更事件。
// Key 采用稳定格式 health-update:{petID}:{status}，使重试与重复调用
// 都能被存储层和下游发布器正确去重。
func NewHealthUpdate(petID, status string) HealthUpdate {
	return HealthUpdate{
		PetID:  petID,
		Status: status,
		Key:    fmt.Sprintf("health-update:%s:%s", petID, status),
	}
}
