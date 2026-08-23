package model

import (
	"context"
	"errors"
)

var ErrRefreshTemporary = errors.New("pet refresh temporarily unavailable")

// ShouldRetryRefresh 判断错误是否值得重试。
// 上下文取消或超时（Canceled / DeadlineExceeded）不可重试：
// 上下文已经失效，重试无意义且会拖慢关闭流程。
func ShouldRetryRefresh(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	return true
}
