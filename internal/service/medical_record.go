package service

import (
	"sort"
	"time"

	"petsmanagement/internal/model"
	"petsmanagement/pkg/idgen"
)

// CreateMedicalRecord 新增就诊记录，校验宠物外键存在。
func (s *Service) CreateMedicalRecord(m model.MedicalRecord) (*model.MedicalRecord, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetPet(m.PetID); err != nil {
		return nil, model.NewValidationError("pet_id", "宠物不存在")
	}
	now := time.Now()
	m.ID = idgen.Hex()
	m.CreatedAt = now
	if err := s.store.CreateMedicalRecord(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMedicalRecord 按 ID 查询就诊记录。
func (s *Service) GetMedicalRecord(id string) (*model.MedicalRecord, error) {
	return s.store.GetMedicalRecord(id)
}

// ListMedicalRecords 按筛选条件查询就诊记录列表，支持分页。
func (s *Service) ListMedicalRecords(filter model.MedicalRecordFilter, page, size int) ([]*model.MedicalRecord, int, error) {
	all := s.store.ListMedicalRecords()
	matched := make([]*model.MedicalRecord, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].VisitAt.After(matched[j].VisitAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.MedicalRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateMedicalRecord 更新就诊记录。
func (s *Service) UpdateMedicalRecord(id string, in model.MedicalRecord) (*model.MedicalRecord, error) {
	existing, err := s.store.GetMedicalRecord(id)
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
	if err := s.store.UpdateMedicalRecord(&in); err != nil {
		return nil, err
	}
	return &in, nil
}

// DeleteMedicalRecord 删除就诊记录。
func (s *Service) DeleteMedicalRecord(id string) error {
	return s.store.DeleteMedicalRecord(id)
}
