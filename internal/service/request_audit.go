package service

import (
	"petsmanagement/internal/adapter"
	"petsmanagement/internal/scopepool"
)

func (s *Service) AuditPetLookup(ownerID, petID string, labels []string, pool *scopepool.Pool, audit *adapter.RequestAudit) {
	scope := pool.Get()
	scope.Prepare(ownerID, petID, labels)
	go audit.Record(scope)
	<-audit.Started()
	pool.Put(scope)
}
