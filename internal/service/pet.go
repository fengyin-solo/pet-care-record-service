package service

import (
	"sort"
	"time"

	"petsmanagement/internal/model"
	"petsmanagement/pkg/idgen"
)

// CreatePet 新增宠物档案，校验主人外键存在。
func (s *Service) CreatePet(p model.Pet) (*model.Pet, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetOwner(p.OwnerID); err != nil {
		return nil, model.NewValidationError("owner_id", "主人不存在")
	}
	now := time.Now()
	p.ID = idgen.Hex()
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := s.store.CreatePet(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetPet 按 ID 查询宠物。
func (s *Service) GetPet(id string) (*model.Pet, error) {
	return s.store.GetPet(id)
}

// ListPets 按筛选条件查询宠物列表，支持分页。
func (s *Service) ListPets(filter model.PetFilter, page, size int) ([]*model.Pet, int, error) {
	all := s.store.ListPets()
	matched := make([]*model.Pet, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Pet{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdatePet 更新宠物档案。
func (s *Service) UpdatePet(id string, in model.Pet) (*model.Pet, error) {
	existing, err := s.store.GetPet(id)
	if err != nil {
		return nil, err
	}
	in.ID = existing.ID
	in.CreatedAt = existing.CreatedAt
	in.UpdatedAt = time.Now()
	if err := in.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetOwner(in.OwnerID); err != nil {
		return nil, model.NewValidationError("owner_id", "主人不存在")
	}
	if err := s.store.UpdatePet(&in); err != nil {
		return nil, err
	}
	return &in, nil
}

// DeletePet 删除宠物档案。
func (s *Service) DeletePet(id string) error {
	return s.store.DeletePet(id)
}

// ArchivePet 归档宠物（active→archived，不可逆）。
func (s *Service) ArchivePet(id string) (*model.Pet, error) {
	existing, err := s.store.GetPet(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionPetStatus(existing.Status, model.PetStatusArchived) {
		return nil, model.NewValidationError("status", "宠物状态不允许归档")
	}
	existing.Status = model.PetStatusArchived
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdatePet(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// BatchArchivePets 批量归档宠物，返回成功数量。
func (s *Service) BatchArchivePets(ids []string) (int, error) {
	success := 0
	for _, id := range ids {
		if _, err := s.ArchivePet(id); err == nil {
			success++
		}
	}
	return success, nil
}
