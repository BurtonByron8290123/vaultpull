package envresolve

import (
	"testing"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv(envAllowFallback, "")
	t.Setenv(envErrorOnMissing, "")
	p := FromEnv()
	if !p.AllowEnvFallback {
		t.Error("expected AllowEnvFallback=true by default")
	}
	if p.ErrorOnMissing {
		t.Error("expected ErrorOnMissing=false by default")
	}
}

func TestFromEnvDisablesFallback(t *testing.T) {
	t.Setenv(envAllowFallback, "false")
	t.Setenv(envErrorOnMissing, "")
	p := FromEnv()
	if p.AllowEnvFallback {
		t.Error("expected AllowEnvFallback=false")
	}
}

func TestFromEnvEnablesErrorOnMissing(t *testing.T) {
	t.Setenv(envAllowFallback, "")
	t.Setenv(envErrorOnMissing, "true")
	p := FromEnv()
	if !p.ErrorOnMissing {
		t.Error("expected ErrorOnMissing=true")
	}
}

func TestFromEnvIgnoresInvalidValues(t *testing.T) {
	t.Setenv(envAllowFallback, "yes")
	t.Setenv(envErrorOnMissing, "1")
	p := FromEnv()
	// neither "yes" nor "1" equals "true", so both should be false
	if p.AllowEnvFallback {
		t.Error("expected AllowEnvFallback=false for invalid value")
	}
	if p.ErrorOnMissing {
		t.Error("expected ErrorOnMissing=false for invalid value")
	}
}
