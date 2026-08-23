// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"petsmanagement/internal/model"
)

var (
	// ErrNotFound 表示记录不存在。
	ErrNotFound = errors.New("记录不存在")
	// ErrConflict 表示记录已存在或状态冲突。
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Owner
	CreateOwner(o *model.Owner) error
	GetOwner(id string) (*model.Owner, error)
	ListOwners() []*model.Owner
	UpdateOwner(o *model.Owner) error
	DeleteOwner(id string) error

	// Pet
	CreatePet(p *model.Pet) error
	GetPet(id string) (*model.Pet, error)
	ListPets() []*model.Pet
	UpdatePet(p *model.Pet) error
	DeletePet(id string) error

	// MedicalRecord
	CreateMedicalRecord(m *model.MedicalRecord) error
	GetMedicalRecord(id string) (*model.MedicalRecord, error)
	ListMedicalRecords() []*model.MedicalRecord
	UpdateMedicalRecord(m *model.MedicalRecord) error
	DeleteMedicalRecord(id string) error

	// Vaccine
	CreateVaccine(v *model.Vaccine) error
	GetVaccine(id string) (*model.Vaccine, error)
	ListVaccines() []*model.Vaccine
	UpdateVaccine(v *model.Vaccine) error
	DeleteVaccine(id string) error

	// FollowUp
	CreateFollowUp(f *model.FollowUp) error
	GetFollowUp(id string) (*model.FollowUp, error)
	ListFollowUps() []*model.FollowUp
	UpdateFollowUp(f *model.FollowUp) error
	DeleteFollowUp(id string) error
	ApplyHealthUpdate(update model.HealthUpdate) bool
	HealthUpdateHistory() []model.HealthUpdate
}
