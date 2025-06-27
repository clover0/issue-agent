package util

import (
	"errors"
	"fmt"
	"time"
)

type RetryableError struct {
	Err      error
	WaitTime time.Duration
}

func (r *RetryableError) Error() string {
	if r.Err != nil {
		return r.Err.Error()
	}
	return "retryable error"
}

func (r *RetryableError) Unwrap() error {
	return r.Err
}

func NewRetryableError(err error, waitTime time.Duration) *RetryableError {
	return &RetryableError{
		Err:      err,
		WaitTime: waitTime,
	}
}

func Retry(maxRetry int, f func() error) error {
	var lastErr error
	for i := 0; i < maxRetry; i++ {
		time.Sleep(time.Duration(i*i) * time.Second)
		if err := f(); err != nil {
			var retryErr *RetryableError
			if errors.As(err, &retryErr) {
				if retryErr.WaitTime > 0 {
					time.Sleep(retryErr.WaitTime)
				}
				lastErr = retryErr
				continue
			}

			return err
		}

		return nil
	}

	return fmt.Errorf("reached maximum retry limit of %d: %w", maxRetry, lastErr)
}
