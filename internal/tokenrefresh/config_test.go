package tokenrefresh

import (
	"testing"
	"time"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv(envThresholdSec, "")
	t.Setenv(envMaxRetries, "")
	p, err := FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.RenewThreshold != 5*time.Minute {
		t.Errorf("expected 5m, got %v", p.RenewThreshold)
	}
	if p.MaxRetries != 3 {
		t.Errorf("expected 3, got %d", p.MaxRetries)
	}
}

func TestFromEnvReadsThreshold(t *testing.T) {
	t.Setenv(envThresholdSec, "120")
	t.Setenv(envMaxRetries, "")
	p, err := FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.RenewThreshold != 120*time.Second {
		t.Errorf("expected 120s, got %v", p.RenewThreshold)
	}
}

func TestFromEnvReadsMaxRetries(t *testing.T) {
	t.Setenv(envThresholdSec, "")
	t.Setenv(envMaxRetries, "5")
	p, err := FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.MaxRetries != 5 {
		t.Errorf("expected 5, got %d", p.MaxRetries)
	}
}

func TestFromEnvInvalidThresholdReturnsError(t *testing.T) {
	t.Setenv(envThresholdSec, "abc")
	_, err := FromEnv()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFromEnvZeroThresholdReturnsError(t *testing.T) {
	t.Setenv(envThresholdSec, "0")
	_, err := FromEnv()
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}
