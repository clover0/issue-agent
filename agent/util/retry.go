package util

import (
	"errors"
	"time"
)

var (
	ErrRetryable          = errors.New("retryable error")
	ErrRetryAfter60Second = errors.New("retry after 60 seconds error")
)

func Retry(maxRetry int, f func() error) error {
	for i := 0; i < maxRetry; i++ {
		time.Sleep(time.Duration(i*i) * time.Second)
		if err := f(); err != nil {
			if errors.Is(err, ErrRetryAfter60Second) {
				time.Sleep(60 * time.Second)
				continue
			}
			if errors.Is(err, ErrRetryable) {
				continue
			}
			return err
		}
		return nil
	}
	return nil
}
