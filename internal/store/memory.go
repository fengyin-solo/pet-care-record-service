package store

import (
	"sync"

	"petsmanagement/internal/model"
)

// MemoryStore 基于内存的 Store 实现，线程安全。
type MemoryStore struct {
	mu             sync.RWMutex
	owners         map[string]*model.Owner
	pets           map[string]*model.Pet
	medicalRecords map[string]*model.MedicalRecord
	vaccines       map[string]*model.Vaccine
	followUps      map[string]*model.FollowUp
}

// NewMemoryStore 创建空的内存 Store。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		owners:         make(map[string]*model.Owner),
		pets:           make(map[string]*model.Pet),
		medicalRecords: make(map[string]*model.MedicalRecord),
		vaccines:       make(map[string]*model.Vaccine),
		followUps:      make(map[string]*model.FollowUp),
	}
}

var _ Store = (*MemoryStore)(nil)
