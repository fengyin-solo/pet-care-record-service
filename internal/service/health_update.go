package service

import (
	"errors"

	"petsmanagement/internal/adapter"
	"petsmanagement/internal/model"
)

// UpdateHealthWithRetry 更新宠物健康状态，发布器繁忙时自动重试一次。
//
// 修复要点：
//   - ApplyHealthUpdate 移到 Publish 成功之后，避免发布失败也写历史。
//   - 发布成功时返回 nil 而非 firstErr，确保接口正确反映最终结果。
//   - 使用稳定标识（petID+status）构造 update，两次尝试共享同一 Key。
func (s *Service) UpdateHealthWithRetry(petID, status string, publisher adapter.Publisher) error {
	update := model.NewHealthUpdate(petID, status)
	var firstErr error
	for i := 0; i < 2; i++ {
		if err := publisher.Publish(update.Key); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			if errors.Is(err, adapter.ErrPublisherBusy) {
				continue
			}
			return err
		}
		s.store.ApplyHealthUpdate(update)
		return nil
	}
	return firstErr
}

func (s *Service) HealthUpdateHistory() []model.HealthUpdate { return s.store.HealthUpdateHistory() }
