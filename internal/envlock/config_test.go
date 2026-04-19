package envlock

import (
	"testing"
	"time"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv(envLockTimeout, "")
	t.Setenv(envLockSuffix, "")

	cfg := FromEnv()

	if cfg.Timeout != defaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", defaultTimeout, cfg.Timeout)
	}
	if cfg.Suffix != defaultLockExt {
		t.Fatalf("expected default suffix %q, got %q", defaultLockExt, cfg.Suffix)
	}
}

func TestFromEnvReadsTimeout(t *testing.T) {
	t.Setenv(envLockTimeout, "30")
	t.Setenv(envLockSuffix, "")

	cfg := FromEnv()

	if cfg.Timeout != 30*time.Second {
		t.Fatalf("expected 30s, got %v", cfg.Timeout)
	}
}

func TestFromEnvReadsSuffix(t *testing.T) {
	t.Setenv(envLockTimeout, "")
	t.Setenv(envLockSuffix, ".lck")

	cfg := FromEnv()

	if cfg.Suffix != ".lck" {
		t.Fatalf("expected .lck, got %q", cfg.Suffix)
	}
}

func TestFromEnvIgnoresInvalidTimeout(t *testing.T) {
	t.Setenv(envLockTimeout, "not-a-number")
	t.Setenv(envLockSuffix, "")

	cfg := FromEnv()

	if cfg.Timeout != defaultTimeout {
		t.Fatalf("expected default timeout on invalid input, got %v", cfg.Timeout)
	}
}

func TestFromEnvIgnoresZeroTimeout(t *testing.T) {
	t.Setenv(envLockTimeout, "0")
	t.Setenv(envLockSuffix, "")

	cfg := FromEnv()

	if cfg.Timeout != defaultTimeout {
		t.Fatalf("expected default timeout for zero value, got %v", cfg.Timeout)
	}
}

func TestDefaultConfigIsValid(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Timeout <= 0 {
		t.Fatal("default timeout must be positive")
	}
	if cfg.Suffix == "" {
		t.Fatal("default suffix must not be empty")
	}
}
