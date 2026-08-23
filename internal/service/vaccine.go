package service

import (
	"sort"
	"time"

	"petsmanagement/internal/model"
	"petsmanagement/pkg/idgen"
)

// CreateVaccine 新增疫苗记录，校验宠物外键存在。
func (s *Service) CreateVaccine(v model.Vaccine) (*model.Vaccine, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetPet(v.PetID); err != nil {
		return nil, model.NewValidationError("pet_id", "宠物不存在")
	}
	now := time.Now()
	v.ID = idgen.Hex()
	v.CreatedAt = now
	if err := s.store.CreateVaccine(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

// GetVaccine 按 ID 查询疫苗记录。
func (s *Service) GetVaccine(id string) (*model.Vaccine, error) {
	return s.store.GetVaccine(id)
}

// ListVaccines 按筛选条件查询疫苗记录列表，支持分页。
func (s *Service) ListVaccines(filter model.VaccineFilter, page, size int) ([]*model.Vaccine, int, error) {
	all := s.store.ListVaccines()
	matched := make([]*model.Vaccine, 0, len(all))
	for _, v := range all {
		if filter.Match(v) {
			matched = append(matched, v)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].AdministeredAt.After(matched[j].AdministeredAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Vaccine{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// ListDueVaccines 查询在指定天数内到期的疫苗，返回到期提醒。
func (s *Service) ListDueVaccines(days int) ([]*model.Vaccine, error) {
	if days <= 0 {
		days = 30
	}
	dueBefore := time.Now().AddDate(0, 0, days)
	items, _, err := s.ListVaccines(model.VaccineFilter{DueBefore: &dueBefore}, 1, 10000)
	return items, err
}

// UpdateVaccine 更新疫苗记录。
func (s *Service) UpdateVaccine(id string, in model.Vaccine) (*model.Vaccine, error) {
	existing, err := s.store.GetVaccine(id)
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
	if err := s.store.UpdateVaccine(&in); err != nil {
		return nil, err
	}
	return &in, nil
}

// DeleteVaccine 删除疫苗记录。
func (s *Service) DeleteVaccine(id string) error {
	return s.store.DeleteVaccine(id)
}
