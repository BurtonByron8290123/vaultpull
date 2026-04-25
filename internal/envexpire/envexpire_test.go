package envexpire

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time { return func() time.Time { return t } }

func newChecker(t *testing.T, p Policy) *Checker {
	t.Helper()
	c, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func TestDefaultPolicyIsValid(t *testing.T) {
	if err := DefaultPolicy().validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStatusOKWhenFresh(t *testing.T) {
	now := time.Now()
	c := newChecker(t, DefaultPolicy())
	c.clock = fixedClock(now)
	results := c.Check(map[string]time.Time{"DB_PASS": now.Add(-1 * 24 * time.Hour)})
	if len(results) != 1 || results[0].Status != StatusOK {
		t.Fatalf("expected OK, got %v", results)
	}
}

func TestStatusWarningWhenApproaching(t *testing.T) {
	now := time.Now()
	p := DefaultPolicy()
	c := newChecker(t, p)
	c.clock = fixedClock(now)
	ts := now.Add(-26 * 24 * time.Hour) // past WarnAge (25d) but before MaxAge (30d)
	results := c.Check(map[string]time.Time{"API_KEY": ts})
	if len(results) != 1 || results[0].Status != StatusWarning {
		t.Fatalf("expected Warning, got %v", results)
	}
}

func TestStatusExpiredWhenOld(t *testing.T) {
	now := time.Now()
	c := newChecker(t, DefaultPolicy())
	c.clock = fixedClock(now)
	ts := now.Add(-31 * 24 * time.Hour)
	results := c.Check(map[string]time.Time{"SECRET": ts})
	if len(results) != 1 || results[0].Status != StatusExpired {
		t.Fatalf("expected Expired, got %v", results)
	}
}

func TestExpiredFiltersResults(t *testing.T) {
	now := time.Now()
	c := newChecker(t, DefaultPolicy())
	c.clock = fixedClock(now)
	timestamps := map[string]time.Time{
		"FRESH":   now.Add(-1 * 24 * time.Hour),
		"STALE":   now.Add(-31 * 24 * time.Hour),
		"WARNING": now.Add(-26 * 24 * time.Hour),
	}
	all := c.Check(timestamps)
	expired := Expired(all)
	if len(expired) != 1 || expired[0].Key != "STALE" {
		t.Fatalf("expected only STALE in expired list, got %v", expired)
	}
}

func TestNewRejectsInvalidPolicy(t *testing.T) {
	_, err := New(Policy{MaxAge: 0, WarnAge: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxAge")
	}
}

func TestValidateRejectsWarnAgeAboveMaxAge(t *testing.T) {
	p := Policy{MaxAge: 10 * time.Hour, WarnAge: 20 * time.Hour}
	if err := p.validate(); err == nil {
		t.Fatal("expected error when WarnAge > MaxAge")
	}
}
