package envrotate

import (
	"errors"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func newRotator(t *testing.T, p Policy) *Rotator {
	t.Helper()
	r, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return r
}

func TestDefaultPolicyIsValid(t *testing.T) {
	_, err := New(DefaultPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStatusOKWhenFresh(t *testing.T) {
	now := time.Now()
	r := newRotator(t, DefaultPolicy())
	r.clock = fixedClock(now)

	results, err := r.Check(map[string]time.Time{
		"DB_PASSWORD": now.Add(-10 * 24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Status != StatusOK {
		t.Errorf("expected StatusOK, got %v", results[0].Status)
	}
}

func TestStatusWarningWhenApproachingMaxAge(t *testing.T) {
	now := time.Now()
	p := DefaultPolicy()
	r := newRotator(t, p)
	r.clock = fixedClock(now)

	results, err := r.Check(map[string]time.Time{
		"API_KEY": now.Add(-80 * 24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Status != StatusWarning {
		t.Errorf("expected StatusWarning, got %v", results[0].Status)
	}
}

func TestStatusExpiredReturnsError(t *testing.T) {
	now := time.Now()
	r := newRotator(t, DefaultPolicy())
	r.clock = fixedClock(now)

	_, err := r.Check(map[string]time.Time{
		"SECRET_KEY": now.Add(-100 * 24 * time.Hour),
	})
	if !errors.Is(err, ErrRotationRequired) {
		t.Errorf("expected ErrRotationRequired, got %v", err)
	}
}

func TestDryRunDoesNotReturnError(t *testing.T) {
	now := time.Now()
	p := DefaultPolicy()
	p.DryRun = true
	r := newRotator(t, p)
	r.clock = fixedClock(now)

	results, err := r.Check(map[string]time.Time{
		"OLD_SECRET": now.Add(-120 * 24 * time.Hour),
	})
	if err != nil {
		t.Errorf("expected no error in dry-run mode, got %v", err)
	}
	if results[0].Status != StatusExpired {
		t.Errorf("expected StatusExpired, got %v", results[0].Status)
	}
}

func TestValidateRejectsWarnAgeAboveMaxAge(t *testing.T) {
	p := Policy{MaxAge: 10 * time.Hour, WarnAge: 20 * time.Hour}
	_, err := New(p)
	if err == nil {
		t.Error("expected error when WarnAge >= MaxAge")
	}
}
