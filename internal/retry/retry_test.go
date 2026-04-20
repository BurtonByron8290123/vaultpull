package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/retry"
)

var errTransient = errors.New("transient")
var errFatal = errors.New("fatal")

func fastPolicy(max int) retry.Policy {
	return retry.Policy{
		MaxAttempts:  max,
		InitialDelay: time.Millisecond,
		MaxDelay:     5 * time.Millisecond,
		Multiplier:   1.5,
	}
}

func TestSucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(3), retry.AlwaysRetry, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetriesOnTransientError(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(3), retry.AlwaysRetry, func() error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestExhaustsMaxAttempts(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(3), retry.AlwaysRetry, func() error {
		calls++
		return errTransient
	})
	if !retry.IsMaxAttempts(err) {
		t.Fatalf("expected MaxAttemptsError, got %v", err)
	}
	if !errors.Is(err, errTransient) {
		t.Fatalf("expected cause to be errTransient, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDoesNotRetryNonRetryableError(t *testing.T) {
	calls := 0
	predicate := func(err error) bool { return !errors.Is(err, errFatal) }
	err := retry.Do(context.Background(), fastPolicy(5), predicate, func() error {
		calls++
		return errFatal
	})
	if !errors.Is(err, errFatal) {
		t.Fatalf("expected errFatal, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestContextCancellationStopsRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := retry.Do(ctx, fastPolicy(5), retry.AlwaysRetry, func() error {
		return errTransient
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestContextDeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	err := retry.Do(ctx, fastPolicy(100), retry.AlwaysRetry, func() error {
		return errTransient
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
}
