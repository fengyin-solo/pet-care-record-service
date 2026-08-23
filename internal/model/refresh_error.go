package model

import "errors"

var ErrRefreshTemporary = errors.New("pet refresh temporarily unavailable")

func ShouldRetryRefresh(err error) bool {
	return err != nil
}
