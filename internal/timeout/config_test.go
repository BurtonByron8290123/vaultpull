package timeout_test

import (
	"testing"
	"time"

	"github.com/vaultpull/vaultpull/internal/timeout"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("VAULTPULL_TIMEOUT_DIAL", "")
	t.Setenv("VAULTPULL_TIMEOUT_READ", "")
	t.Setenv("VAULTPULL_TIMEOUT_OVERALL", "")

	p, err := timeout.FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	def := timeout.DefaultPolicy()
	if p.Dial != def.Dial || p.Read != def.Read || p.Overall != def.Overall {
		t.Fatalf("expected default policy, got %+v", p)
	}
}

func TestFromEnvReadsDial(t *testing.T) {
	t.Setenv("VAULTPULL_TIMEOUT_DIAL", "3")
	t.Setenv("VAULTPULL_TIMEOUT_READ", "")
	t.Setenv("VAULTPULL_TIMEOUT_OVERALL", "")

	p, err := timeout.FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Dial != 3*time.Second {
		t.Fatalf("expected 3s dial, got %v", p.Dial)
	}
}

func TestFromEnvReadsOverall(t *testing.T) {
	t.Setenv("VAULTPULL_TIMEOUT_OVERALL", "60")
	t.Setenv("VAULTPULL_TIMEOUT_DIAL", "")
	t.Setenv("VAULTPULL_TIMEOUT_READ", "")

	p, err := timeout.FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Overall != 60*time.Second {
		t.Fatalf("expected 60s overall, got %v", p.Overall)
	}
}

func TestFromEnvInvalidValueReturnsError(t *testing.T) {
	t.Setenv("VAULTPULL_TIMEOUT_DIAL", "notanumber")
	t.Setenv("VAULTPULL_TIMEOUT_READ", "")
	t.Setenv("VAULTPULL_TIMEOUT_OVERALL", "")

	_, err := timeout.FromEnv()
	if err == nil {
		t.Fatal("expected error for invalid dial value")
	}
}
