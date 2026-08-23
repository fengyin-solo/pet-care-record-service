// Package service 实现业务逻辑层。
package service

import (
	"petsmanagement/internal/config"
	"petsmanagement/internal/store"
	"petsmanagement/pkg/logger"
)

// Service 业务服务，聚合各实体业务方法。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 创建业务服务。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
