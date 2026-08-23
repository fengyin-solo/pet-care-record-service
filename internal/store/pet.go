package store

import "petsmanagement/internal/model"

// CreatePet 新增宠物档案。
func (s *MemoryStore) CreatePet(p *model.Pet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pets[p.ID] = p
	return nil
}

// GetPet 按 ID 查询宠物。
func (s *MemoryStore) GetPet(id string) (*model.Pet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.pets[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

// ListPets 返回全部宠物。
func (s *MemoryStore) ListPets() []*model.Pet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Pet, 0, len(s.pets))
	for _, p := range s.pets {
		list = append(list, p)
	}
	return list
}

// UpdatePet 更新宠物档案。
func (s *MemoryStore) UpdatePet(p *model.Pet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.pets[p.ID]; !ok {
		return ErrNotFound
	}
	s.pets[p.ID] = p
	return nil
}

// DeletePet 删除宠物档案。
func (s *MemoryStore) DeletePet(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.pets[id]; !ok {
		return ErrNotFound
	}
	delete(s.pets, id)
	return nil
}
