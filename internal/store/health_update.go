package store

import (
	"fmt"

	"petsmanagement/internal/model"
)

// ApplyHealthUpdate 将健康状态变更写入历史，基于 petID+status 做幂等去重。
// 同一宠物和状态的重复记录（包括 legacy-attempt-N 格式的 key）一律拒绝。
func (s *MemoryStore) ApplyHealthUpdate(update model.HealthUpdate) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := fmt.Sprintf("%s:%s", update.PetID, update.Status)
	if s.healthUpdateKeys[id] {
		return false
	}
	s.healthUpdateKeys[id] = true
	s.healthUpdateHistory = append(s.healthUpdateHistory, update)
	return true
}

func (s *MemoryStore) HealthUpdateHistory() []model.HealthUpdate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.HealthUpdate(nil), s.healthUpdateHistory...)
}
