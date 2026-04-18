package backoff_test

import (
	"testing"
	"time"

	"github.com/vaultpull/vaultpull/internal/backoff"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	p := backoff.FromEnv()
	def := backoff.DefaultPolicy()
	if p.InitialInterval != def.InitialInterval {
		t.Fatalf("expected default initial interval, got %v", p.InitialInterval)
	}
}

func TestFromEnvReadsInitialMS(t *testing.T) {
	t.Setenv("VAULTPULL_BACKOFF_INITIAL_MS", "500")
	p := backoff.FromEnv()
	if p.InitialInterval != 500*time.Millisecond {
		t.Fatalf("expected 500ms, got %v", p.InitialInterval)
	}
}

func TestFromEnvReadsMultiplier(t *testing.T) {
	t.Setenv("VAULTPULL_BACKOFF_MULTIPLIER", "3")
	p := backoff.FromEnv()
	if p.Multiplier != 3.0 {
		t.Fatalf("expected multiplier 3, got %v", p.Multiplier)
	}
}

func TestFromEnvIgnoresInvalidMultiplier(t *testing.T) {
	t.Setenv("VAULTPULL_BACKOFF_MULTIPLIER", "0.1") // below 1, should be ignored
	p := backoff.FromEnv()
	if p.Multiplier != backoff.DefaultPolicy().Multiplier {
		t.Fatalf("expected default multiplier, got %v", p.Multiplier)
	}
}

func TestFromEnvReadsMaxMS(t *testing.T) {
	t.Setenv("VAULTPULL_BACKOFF_MAX_MS", "60000")
	p := backoff.FromEnv()
	if p.MaxInterval != 60*time.Second {
		t.Fatalf("expected 60s, got %v", p.MaxInterval)
	}
}

func TestFromEnvReadsJitter(t *testing.T) {
	t.Setenv("VAULTPULL_BACKOFF_JITTER", "0.5")
	p := backoff.FromEnv()
	if p.Jitter != 0.5 {
		t.Fatalf("expected jitter 0.5, got %v", p.Jitter)
	}
}
