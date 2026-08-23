package store

import "petsmanagement/internal/model"

// CreateMedicalRecord 新增就诊记录。
func (s *MemoryStore) CreateMedicalRecord(m *model.MedicalRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.medicalRecords[m.ID] = m
	return nil
}

// GetMedicalRecord 按 ID 查询就诊记录。
func (s *MemoryStore) GetMedicalRecord(id string) (*model.MedicalRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.medicalRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

// ListMedicalRecords 返回全部就诊记录。
func (s *MemoryStore) ListMedicalRecords() []*model.MedicalRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.MedicalRecord, 0, len(s.medicalRecords))
	for _, m := range s.medicalRecords {
		list = append(list, m)
	}
	return list
}

// UpdateMedicalRecord 更新就诊记录。
func (s *MemoryStore) UpdateMedicalRecord(m *model.MedicalRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.medicalRecords[m.ID]; !ok {
		return ErrNotFound
	}
	s.medicalRecords[m.ID] = m
	return nil
}

// DeleteMedicalRecord 删除就诊记录。
func (s *MemoryStore) DeleteMedicalRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.medicalRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.medicalRecords, id)
	return nil
}
