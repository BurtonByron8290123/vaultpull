package retry

import (
	"context"
	"errors"
	"time"
)

// Policy defines the retry behaviour for transient errors.
type Policy struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts:  4,
		InitialDelay: 250 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// IsRetryable is a predicate that callers can provide to decide whether a
// given error should trigger a retry.
type IsRetryable func(err error) bool

// AlwaysRetry is an IsRetryable that retries on every non-nil error.
func AlwaysRetry(err error) bool { return err != nil }

// Do executes fn according to p, retrying on errors deemed retryable by the
// supplied predicate. The context is checked before every attempt.
func Do(ctx context.Context, p Policy, retryable IsRetryable, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	delay := p.InitialDelay
	var lastErr error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if !retryable(lastErr) {
			return lastErr
		}
		if attempt == p.MaxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		delay = time.Duration(float64(delay) * p.Multiplier)
		if delay > p.MaxDelay {
			delay = p.MaxDelay
		}
	}
	return &MaxAttemptsError{Attempts: p.MaxAttempts, Cause: lastErr}
}

// MaxAttemptsError is returned when all retry attempts have been exhausted.
type MaxAttemptsError struct {
	Attempts int
	Cause    error
}

func (e *MaxAttemptsError) Error() string {
	return "retry: max attempts reached: " + e.Cause.Error()
}

func (e *MaxAttemptsError) Unwrap() error { return e.Cause }

// IsMaxAttempts reports whether err is a MaxAttemptsError.
func IsMaxAttempts(err error) bool {
	var t *MaxAttemptsError
	return errors.As(err, &t)
}
