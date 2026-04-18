package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestDefaultPolicyIsValid(t *testing.T) {
	p := DefaultPolicy()
	if err := p.Validate(); err != nil {
		t.Fatalf("expected valid default policy, got: %v", err)
	}
}

func TestValidateRejectsZeroRPS(t *testing.T) {
	p := Policy{RequestsPerSecond: 0, Burst: 1}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero RPS")
	}
}

func TestValidateRejectsZeroBurst(t *testing.T) {
	p := Policy{RequestsPerSecond: 1, Burst: 0}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero burst")
	}
}

func TestNewReturnsErrorOnInvalidPolicy(t *testing.T) {
	_, err := New(Policy{RequestsPerSecond: -1, Burst: 1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWaitConsumesToken(t *testing.T) {
	l, err := New(Policy{RequestsPerSecond: 100, Burst: 2})
	if err != nil {
		t.Fatal(err)
	}
	// Replace waitFunc so tests never sleep.
	l.waitFunc = func(_ context.Context, _ time.Duration) error { return nil }

	for i := 0; i < 5; i++ {
		if err := l.Wait(context.Background()); err != nil {
			t.Fatalf("Wait() returned unexpected error on call %d: %v", i, err)
		}
	}
}

func TestWaitRespectsContextCancellation(t *testing.T) {
	l, err := New(Policy{RequestsPerSecond: 0.001, Burst: 1})
	if err != nil {
		t.Fatal(err)
	}
	// Drain the burst token.
	_ = l.Wait(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := l.Wait(ctx); err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestTokensRefillOverTime(t *testing.T) {
	l, err := New(Policy{RequestsPerSecond: 1000, Burst: 1})
	if err != nil {
		t.Fatal(err)
	}
	l.waitFunc = func(_ context.Context, _ time.Duration) error { return nil }
	// Drain burst.
	_ = l.Wait(context.Background())
	// Simulate time passing by rewinding last.
	l.last = time.Now().Add(-10 * time.Millisecond)
	if err := l.Wait(context.Background()); err != nil {
		t.Fatal("expected token to be available after refill")
	}
}
