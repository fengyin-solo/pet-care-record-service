package store

import "petsmanagement/internal/model"

// CreateVaccine 新增疫苗记录。
func (s *MemoryStore) CreateVaccine(v *model.Vaccine) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vaccines[v.ID] = v
	return nil
}

// GetVaccine 按 ID 查询疫苗记录。
func (s *MemoryStore) GetVaccine(id string) (*model.Vaccine, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.vaccines[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

// ListVaccines 返回全部疫苗记录。
func (s *MemoryStore) ListVaccines() []*model.Vaccine {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Vaccine, 0, len(s.vaccines))
	for _, v := range s.vaccines {
		list = append(list, v)
	}
	return list
}

// UpdateVaccine 更新疫苗记录。
func (s *MemoryStore) UpdateVaccine(v *model.Vaccine) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.vaccines[v.ID]; !ok {
		return ErrNotFound
	}
	s.vaccines[v.ID] = v
	return nil
}

// DeleteVaccine 删除疫苗记录。
func (s *MemoryStore) DeleteVaccine(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.vaccines[id]; !ok {
		return ErrNotFound
	}
	delete(s.vaccines, id)
	return nil
}
