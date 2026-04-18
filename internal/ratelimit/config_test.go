package ratelimit

import (
	"testing"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv(envRPS, "")
	t.Setenv(envBurst, "")

	p, err := FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	def := DefaultPolicy()
	if p.RequestsPerSecond != def.RequestsPerSecond {
		t.Errorf("RPS: got %f, want %f", p.RequestsPerSecond, def.RequestsPerSecond)
	}
	if p.Burst != def.Burst {
		t.Errorf("Burst: got %d, want %d", p.Burst, def.Burst)
	}
}

func TestFromEnvReadsRPS(t *testing.T) {
	t.Setenv(envRPS, "25.5")
	t.Setenv(envBurst, "")

	p, err := FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.RequestsPerSecond != 25.5 {
		t.Errorf("got %f, want 25.5", p.RequestsPerSecond)
	}
}

func TestFromEnvReadsBurst(t *testing.T) {
	t.Setenv(envRPS, "")
	t.Setenv(envBurst, "20")

	p, err := FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Burst != 20 {
		t.Errorf("got %d, want 20", p.Burst)
	}
}

func TestFromEnvInvalidRPSReturnsError(t *testing.T) {
	t.Setenv(envRPS, "not-a-number")
	if _, err := FromEnv(); err == nil {
		t.Fatal("expected error for invalid RPS")
	}
}

func TestFromEnvInvalidBurstReturnsError(t *testing.T) {
	t.Setenv(envRPS, "")
	t.Setenv(envBurst, "bad")
	if _, err := FromEnv(); err == nil {
		t.Fatal("expected error for invalid burst")
	}
}

func TestFromEnvZeroRPSFailsValidation(t *testing.T) {
	t.Setenv(envRPS, "0")
	t.Setenv(envBurst, "")
	if _, err := FromEnv(); err == nil {
		t.Fatal("expected validation error for zero RPS")
	}
}
