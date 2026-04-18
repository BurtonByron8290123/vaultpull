package expire_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/expire"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestDefaultPolicyIsValid(t *testing.T) {
	p := expire.DefaultPolicy()
	if err := p.Validate(); err != nil {
		t.Fatalf("expected valid policy, got %v", err)
	}
}

func TestIsExpiredReturnsFalseWhenFresh(t *testing.T) {
	now := time.Now()
	p := expire.Policy{TTL: time.Minute, ClockFunc: fixedClock(now)}
	fetchedAt := now.Add(-30 * time.Second)
	if p.IsExpired(fetchedAt) {
		t.Fatal("expected not expired")
	}
}

func TestIsExpiredReturnsTrueWhenStale(t *testing.T) {
	now := time.Now()
	p := expire.Policy{TTL: time.Minute, ClockFunc: fixedClock(now)}
	fetchedAt := now.Add(-2 * time.Minute)
	if !p.IsExpired(fetchedAt) {
		t.Fatal("expected expired")
	}
}

func TestExpiresAtIsCorrect(t *testing.T) {
	now := time.Now()
	p := expire.Policy{TTL: time.Minute, ClockFunc: fixedClock(now)}
	fetchedAt := now
	want := fetchedAt.Add(time.Minute)
	if got := p.ExpiresAt(fetchedAt); !got.Equal(want) {
		t.Fatalf("ExpiresAt: got %v, want %v", got, want)
	}
}

func TestValidateRejectsZeroTTL(t *testing.T) {
	p := expire.Policy{TTL: 0, ClockFunc: time.Now}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestValidateRejectsNilClock(t *testing.T) {
	p := expire.Policy{TTL: time.Minute, ClockFunc: nil}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for nil clock")
	}
}

func TestFromEnvUsesDefault(t *testing.T) {
	t.Setenv("VAULTPULL_SECRET_TTL_SECONDS", "")
	p, err := expire.FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.TTL != 5*time.Minute {
		t.Fatalf("expected 5m TTL, got %v", p.TTL)
	}
}

func TestFromEnvReadsCustomTTL(t *testing.T) {
	t.Setenv("VAULTPULL_SECRET_TTL_SECONDS", "120")
	p, err := expire.FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.TTL != 2*time.Minute {
		t.Fatalf("expected 2m TTL, got %v", p.TTL)
	}
}

func TestFromEnvInvalidValueReturnsError(t *testing.T) {
	t.Setenv("VAULTPULL_SECRET_TTL_SECONDS", "not-a-number")
	_, err := expire.FromEnv()
	if err == nil {
		t.Fatal("expected error for invalid TTL")
	}
}
