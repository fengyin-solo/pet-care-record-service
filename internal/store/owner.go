package store

import "petsmanagement/internal/model"

// CreateOwner 新增主人。
func (s *MemoryStore) CreateOwner(o *model.Owner) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.owners[o.ID] = o
	return nil
}

// GetOwner 按 ID 查询主人。
func (s *MemoryStore) GetOwner(id string) (*model.Owner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.owners[id]
	if !ok {
		return nil, ErrNotFound
	}
	return o, nil
}

// ListOwners 返回全部主人。
func (s *MemoryStore) ListOwners() []*model.Owner {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Owner, 0, len(s.owners))
	for _, o := range s.owners {
		list = append(list, o)
	}
	return list
}

// UpdateOwner 更新主人。
func (s *MemoryStore) UpdateOwner(o *model.Owner) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.owners[o.ID]; !ok {
		return ErrNotFound
	}
	s.owners[o.ID] = o
	return nil
}

// DeleteOwner 删除主人。
func (s *MemoryStore) DeleteOwner(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.owners[id]; !ok {
		return ErrNotFound
	}
	delete(s.owners, id)
	return nil
}
