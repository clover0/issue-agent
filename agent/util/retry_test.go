package util

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/clover0/issue-agent/test/assert"
)

func TestRetryableError_Error(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err      error
		waitTime time.Duration
		want     string
	}{
		"with underlying error": {
			err:      errors.New("test error"),
			waitTime: time.Second,
			want:     "test error",
		},
		"without underlying error": {
			err:      nil,
			waitTime: time.Second,
			want:     "retryable error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			retryErr := RetryableError{
				Err:      tt.err,
				WaitTime: tt.waitTime,
			}

			got := retryErr.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRetryableError_Unwrap(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("original error")
	retryErr := RetryableError{
		Err:      originalErr,
		WaitTime: time.Second,
	}

	unwrappedErr := retryErr.Unwrap()

	assert.Equal(t, originalErr, unwrappedErr)
}

func TestNewRetryableError(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("original error")
	waitTime := 2 * time.Second

	retryErr := NewRetryableError(originalErr, waitTime)

	assert.Equal(t, originalErr, retryErr.Err)
	assert.Equal(t, waitTime, retryErr.WaitTime)
}

func TestRetry(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		maxRetry      int
		errorSequence []error
		wantCalls     int
		wantResult    error
	}{
		"success on first try": {
			maxRetry:      3,
			errorSequence: []error{nil},
			wantCalls:     1,
		},
		"success after retries": {
			maxRetry:      3,
			errorSequence: []error{NewRetryableError(errors.New("retry"), 10*time.Millisecond), NewRetryableError(errors.New("retry"), 10*time.Millisecond), nil},
			wantCalls:     3,
		},
		"non-retryable error": {
			maxRetry:      3,
			errorSequence: []error{errors.New("non-retryable")},
			wantCalls:     1,
			wantResult:    errors.New("non-retryable"),
		},
		"retryable error with success": {
			maxRetry: 3,
			errorSequence: []error{
				NewRetryableError(errors.New("retry1"), 10*time.Millisecond),
				NewRetryableError(errors.New("retry2"), 10*time.Millisecond),
				nil,
			},
			wantCalls: 3,
		},
		"max retries exceeded": {
			maxRetry: 2,
			errorSequence: []error{
				NewRetryableError(errors.New("retry1"), 10*time.Millisecond),
				NewRetryableError(errors.New("retry2"), 10*time.Millisecond),
				NewRetryableError(errors.New("retry3"), 10*time.Millisecond),
			},
			wantCalls:  2,
			wantResult: fmt.Errorf("reached maximum retry limit of 2: %w", NewRetryableError(errors.New("retry2"), 10*time.Millisecond)),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			callCount := 0
			f := func() error {
				if callCount < len(tt.errorSequence) {
					err := tt.errorSequence[callCount]
					callCount++
					return err
				}
				return nil
			}

			result := Retry(tt.maxRetry, f)

			assert.Equal(t, tt.wantCalls, callCount)

			if tt.wantResult != nil {
				assert.Equal(t, tt.wantResult.Error(), result.Error())
				return
			}

			assert.Nil(t, result)

		})
	}
}
