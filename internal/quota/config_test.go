package quota_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/quota"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("VAULTPULL_QUOTA_MAX_REQUESTS", "")
	p, err := quota.FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if p.MaxRequestsPerPath != 10 {
		t.Fatalf("expected default 10, got %d", p.MaxRequestsPerPath)
	}
}

func TestFromEnvReadsMaxRequests(t *testing.T) {
	t.Setenv("VAULTPULL_QUOTA_MAX_REQUESTS", "25")
	p, err := quota.FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if p.MaxRequestsPerPath != 25 {
		t.Fatalf("expected 25, got %d", p.MaxRequestsPerPath)
	}
}

func TestFromEnvInvalidValueReturnsError(t *testing.T) {
	t.Setenv("VAULTPULL_QUOTA_MAX_REQUESTS", "banana")
	_, err := quota.FromEnv()
	if err == nil {
		t.Fatal("expected error for invalid value")
	}
}

func TestFromEnvZeroValueReturnsError(t *testing.T) {
	t.Setenv("VAULTPULL_QUOTA_MAX_REQUESTS", "0")
	_, err := quota.FromEnv()
	if err == nil {
		t.Fatal("expected error for zero value")
	}
}
