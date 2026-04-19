package leasetrack

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time { return func() time.Time { return t } }

func freshEntry(now time.Time, ttl time.Duration) Entry {
	return Entry{
		Path:     "secret/app",
		LeaseID:  "lease-1",
		IssuedAt: now,
		LeaseTTL: ttl,
	}
}

func TestCheckFreshLease(t *testing.T) {
	now := time.Now()
	tr, _ := New(DefaultPolicy())
	tr.clock = fixedClock(now)
	e := freshEntry(now, 30*time.Minute)
	if got := tr.Check(e); got != StatusFresh {
		t.Fatalf("expected Fresh, got %d", got)
	}
}

func TestCheckWarningLease(t *testing.T) {
	now := time.Now()
	tr, _ := New(DefaultPolicy())
	// Move clock to 5 minutes before expiry (warn threshold is 10 min)
	tr.clock = fixedClock(now.Add(25 * time.Minute))
	e := freshEntry(now, 30*time.Minute)
	if got := tr.Check(e); got != StatusWarning {
		t.Fatalf("expected Warning, got %d", got)
	}
}

func TestCheckExpiredLease(t *testing.T) {
	now := time.Now()
	tr, _ := New(DefaultPolicy())
	tr.clock = fixedClock(now.Add(31 * time.Minute))
	e := freshEntry(now, 30*time.Minute)
	if got := tr.Check(e); got != StatusExpired {
		t.Fatalf("expected Expired, got %d", got)
	}
}

func TestCheckAllFiltersOnlyNonFresh(t *testing.T) {
	now := time.Now()
	tr, _ := New(DefaultPolicy())
	tr.clock = fixedClock(now.Add(26 * time.Minute))
	entries := []Entry{
		freshEntry(now, 30*time.Minute),  // warning
		freshEntry(now, 60*time.Minute),  // fresh
		freshEntry(now, 25*time.Minute),  // expired
	}
	out := tr.CheckAll(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 non-fresh entries, got %d", len(out))
	}
}

func TestNewRejectsInvalidPolicy(t *testing.T) {
	_, err := New(Policy{WarnThreshold: 0})
	if err == nil {
		t.Fatal("expected error for zero WarnThreshold")
	}
}

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv(envWarnThresholdSec, "")
	p := FromEnv()
	if p.WarnThreshold != 10*time.Minute {
		t.Fatalf("unexpected default: %v", p.WarnThreshold)
	}
}

func TestFromEnvReadsWarnThreshold(t *testing.T) {
	t.Setenv(envWarnThresholdSec, "120")
	p := FromEnv()
	if p.WarnThreshold != 120*time.Second {
		t.Fatalf("unexpected threshold: %v", p.WarnThreshold)
	}
}
