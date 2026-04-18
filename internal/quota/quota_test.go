package quota_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpull/internal/quota"
)

func TestAllowUnderLimit(t *testing.T) {
	tr, err := quota.New(quota.Policy{MaxRequestsPerPath: 3})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if err := tr.Allow("secret/app"); err != nil {
			t.Fatalf("unexpected error on attempt %d: %v", i+1, err)
		}
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	tr, _ := quota.New(quota.Policy{MaxRequestsPerPath: 2})
	tr.Allow("secret/app") // 1
	tr.Allow("secret/app") // 2
	err := tr.Allow("secret/app") // 3 — over limit
	if !errors.Is(err, quota.ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestAllowIndependentPaths(t *testing.T) {
	tr, _ := quota.New(quota.Policy{MaxRequestsPerPath: 1})
	if err := tr.Allow("secret/a"); err != nil {
		t.Fatal(err)
	}
	if err := tr.Allow("secret/b"); err != nil {
		t.Fatal(err)
	}
}

func TestCountReturnsAccurateValue(t *testing.T) {
	tr, _ := quota.New(quota.DefaultPolicy())
	tr.Allow("secret/x")
	tr.Allow("secret/x")
	if got := tr.Count("secret/x"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestResetClearsCounters(t *testing.T) {
	tr, _ := quota.New(quota.DefaultPolicy())
	tr.Allow("secret/x")
	tr.Reset()
	if got := tr.Count("secret/x"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestValidateRejectsZeroMax(t *testing.T) {
	_, err := quota.New(quota.Policy{MaxRequestsPerPath: 0})
	if err == nil {
		t.Fatal("expected error for zero MaxRequestsPerPath")
	}
}
