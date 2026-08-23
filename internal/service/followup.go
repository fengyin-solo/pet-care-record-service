package service

import (
	"sort"
	"time"

	"petsmanagement/internal/model"
	"petsmanagement/pkg/idgen"
)

// CreateFollowUp 新增回访记录，校验宠物外键存在。
func (s *Service) CreateFollowUp(f model.FollowUp) (*model.FollowUp, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetPet(f.PetID); err != nil {
		return nil, model.NewValidationError("pet_id", "宠物不存在")
	}
	now := time.Now()
	f.ID = idgen.Hex()
	f.CreatedAt = now
	if err := s.store.CreateFollowUp(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

// GetFollowUp 按 ID 查询回访记录。
func (s *Service) GetFollowUp(id string) (*model.FollowUp, error) {
	return s.store.GetFollowUp(id)
}

// ListFollowUps 按筛选条件查询回访记录列表，支持分页。
func (s *Service) ListFollowUps(filter model.FollowUpFilter, page, size int) ([]*model.FollowUp, int, error) {
	all := s.store.ListFollowUps()
	matched := make([]*model.FollowUp, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].ScheduledAt.Before(matched[j].ScheduledAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.FollowUp{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateFollowUp 更新回访记录。
func (s *Service) UpdateFollowUp(id string, in model.FollowUp) (*model.FollowUp, error) {
	existing, err := s.store.GetFollowUp(id)
	if err != nil {
		return nil, err
	}
	in.ID = existing.ID
	in.CreatedAt = existing.CreatedAt
	if err := in.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetPet(in.PetID); err != nil {
		return nil, model.NewValidationError("pet_id", "宠物不存在")
	}
	if err := s.store.UpdateFollowUp(&in); err != nil {
		return nil, err
	}
	return &in, nil
}

// DeleteFollowUp 删除回访记录。
func (s *Service) DeleteFollowUp(id string) error {
	return s.store.DeleteFollowUp(id)
}

// TransitionFollowUp 流转回访状态（pending→done/cancelled）。
func (s *Service) TransitionFollowUp(id, target string) (*model.FollowUp, error) {
	existing, err := s.store.GetFollowUp(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionFollowUpStatus(existing.Status, target) {
		return nil, model.NewValidationError("status", "回访状态不允许该流转")
	}
	existing.Status = target
	if target == model.FollowUpDone {
		now := time.Now()
		existing.CompletedAt = &now
	}
	if err := s.store.UpdateFollowUp(existing); err != nil {
		return nil, err
	}
	return existing, nil
}
