package backoff_test

import (
	"context"
	"testing"
	"time"

	"github.com/vaultpull/vaultpull/internal/backoff"
)

func TestDefaultPolicyIsValid(t *testing.T) {
	if err := backoff.DefaultPolicy().Validate(); err != nil {
		t.Fatalf("default policy invalid: %v", err)
	}
}

func TestValidateRejectsZeroInitial(t *testing.T) {
	p := backoff.DefaultPolicy()
	p.InitialInterval = 0
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero initial interval")
	}
}

func TestValidateRejectsMultiplierBelowOne(t *testing.T) {
	p := backoff.DefaultPolicy()
	p.Multiplier = 0.5
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for multiplier < 1")
	}
}

func TestValidateRejectsMaxBelowInitial(t *testing.T) {
	p := backoff.DefaultPolicy()
	p.MaxInterval = p.InitialInterval / 2
	if err := p.Validate(); err == nil {
		t.Fatal("expected error when max < initial")
	}
}

func TestNextGrowsWithAttempt(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 100 * time.Millisecond,
		Multiplier:      2.0,
		MaxInterval:     10 * time.Second,
		Jitter:          0, // deterministic
	}
	prev := p.Next(0)
	for i := 1; i <= 4; i++ {
		cur := p.Next(i)
		if cur <= prev {
			t.Fatalf("attempt %d: expected growth, got %v <= %v", i, cur, prev)
		}
		prev = cur
	}
}

func TestNextCapsAtMaxInterval(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		Multiplier:      10.0,
		MaxInterval:     2 * time.Second,
		Jitter:          0,
	}
	for i := 0; i < 10; i++ {
		if d := p.Next(i); d > p.MaxInterval {
			t.Fatalf("attempt %d: %v exceeds max %v", i, d, p.MaxInterval)
		}
	}
}

func TestSleepRespectsContextCancellation(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 10 * time.Second,
		Multiplier:      1.0,
		MaxInterval:     10 * time.Second,
		Jitter:          0,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := p.Sleep(ctx, 0)
	if err == nil {
		t.Fatal("expected context error")
	}
	if time.Since(start) >= 5*time.Second {
		t.Fatal("sleep did not respect cancellation")
	}
}
