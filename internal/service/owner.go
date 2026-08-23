package service

import (
	"sort"
	"time"

	"petsmanagement/internal/model"
	"petsmanagement/pkg/idgen"
)

// CreateOwner 新增主人。
func (s *Service) CreateOwner(o model.Owner) (*model.Owner, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	o.ID = idgen.Hex()
	o.CreatedAt = now
	o.UpdatedAt = now
	if err := s.store.CreateOwner(&o); err != nil {
		return nil, err
	}
	return &o, nil
}

// GetOwner 按 ID 查询主人。
func (s *Service) GetOwner(id string) (*model.Owner, error) {
	return s.store.GetOwner(id)
}

// ListOwners 按筛选条件查询主人列表，支持分页。
func (s *Service) ListOwners(filter model.OwnerFilter, page, size int) ([]*model.Owner, int, error) {
	all := s.store.ListOwners()
	matched := make([]*model.Owner, 0, len(all))
	for _, o := range all {
		if filter.Match(o) {
			matched = append(matched, o)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Owner{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateOwner 更新主人。
func (s *Service) UpdateOwner(id string, in model.Owner) (*model.Owner, error) {
	existing, err := s.store.GetOwner(id)
	if err != nil {
		return nil, err
	}
	in.ID = existing.ID
	in.CreatedAt = existing.CreatedAt
	in.UpdatedAt = time.Now()
	if err := in.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateOwner(&in); err != nil {
		return nil, err
	}
	return &in, nil
}

// DeleteOwner 删除主人。
func (s *Service) DeleteOwner(id string) error {
	return s.store.DeleteOwner(id)
}
