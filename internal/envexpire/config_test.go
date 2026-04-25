package envexpire

import (
	"testing"
	"time"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRE_MAX_AGE_DAYS", "")
	t.Setenv("VAULTPULL_EXPIRE_WARN_AGE_DAYS", "")
	p := FromEnv()
	def := DefaultPolicy()
	if p.MaxAge != def.MaxAge || p.WarnAge != def.WarnAge {
		t.Fatalf("expected defaults, got max=%v warn=%v", p.MaxAge, p.WarnAge)
	}
}

func TestFromEnvReadsMaxAgeDays(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRE_MAX_AGE_DAYS", "60")
	t.Setenv("VAULTPULL_EXPIRE_WARN_AGE_DAYS", "")
	p := FromEnv()
	if p.MaxAge != 60*24*time.Hour {
		t.Fatalf("expected 60 days, got %v", p.MaxAge)
	}
}

func TestFromEnvReadsWarnAgeDays(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRE_MAX_AGE_DAYS", "60")
	t.Setenv("VAULTPULL_EXPIRE_WARN_AGE_DAYS", "50")
	p := FromEnv()
	if p.WarnAge != 50*24*time.Hour {
		t.Fatalf("expected 50 days, got %v", p.WarnAge)
	}
}

func TestFromEnvIgnoresInvalidMaxAge(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRE_MAX_AGE_DAYS", "notanumber")
	t.Setenv("VAULTPULL_EXPIRE_WARN_AGE_DAYS", "")
	p := FromEnv()
	if p.MaxAge != DefaultPolicy().MaxAge {
		t.Fatalf("expected default MaxAge, got %v", p.MaxAge)
	}
}

func TestFromEnvClampsWarnAgeWhenAboveMaxAge(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRE_MAX_AGE_DAYS", "10")
	t.Setenv("VAULTPULL_EXPIRE_WARN_AGE_DAYS", "20")
	p := FromEnv()
	// WarnAge should fall back to default because 20 > 10
	if p.WarnAge >= p.MaxAge {
		t.Fatalf("WarnAge (%v) must be less than MaxAge (%v)", p.WarnAge, p.MaxAge)
	}
}
