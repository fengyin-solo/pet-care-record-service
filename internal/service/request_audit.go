package service

import (
	"petsmanagement/internal/adapter"
	"petsmanagement/internal/scopepool"
)

// AuditPetLookup 审计一次宠物查询。
//
// 从池中取出 scope 并通过 Prepare 交接所有权，启动 Record goroutine
// 后等待 Started 信号。由于 Record 在阻塞前已拷贝身份快照，
// <-Started 之后即可安全 Put 回收 scope，不会影响本次审计内容。
func (s *Service) AuditPetLookup(ownerID, petID string, labels []string, pool *scopepool.Pool, audit *adapter.RequestAudit) {
	scope := pool.Get()
	scope.Prepare(ownerID, petID, labels)
	go audit.Record(scope)
	<-audit.Started()
	pool.Put(scope)
}
