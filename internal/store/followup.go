package store

import "petsmanagement/internal/model"

// CreateFollowUp 新增回访记录。
func (s *MemoryStore) CreateFollowUp(f *model.FollowUp) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.followUps[f.ID] = f
	return nil
}

// GetFollowUp 按 ID 查询回访记录。
func (s *MemoryStore) GetFollowUp(id string) (*model.FollowUp, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.followUps[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

// ListFollowUps 返回全部回访记录。
func (s *MemoryStore) ListFollowUps() []*model.FollowUp {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.FollowUp, 0, len(s.followUps))
	for _, f := range s.followUps {
		list = append(list, f)
	}
	return list
}

// UpdateFollowUp 更新回访记录。
func (s *MemoryStore) UpdateFollowUp(f *model.FollowUp) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.followUps[f.ID]; !ok {
		return ErrNotFound
	}
	s.followUps[f.ID] = f
	return nil
}

// DeleteFollowUp 删除回访记录。
func (s *MemoryStore) DeleteFollowUp(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.followUps[id]; !ok {
		return ErrNotFound
	}
	delete(s.followUps, id)
	return nil
}
