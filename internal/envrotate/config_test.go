package envrotate

import (
	"testing"
	"time"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("VAULTPULL_ROTATE_MAX_AGE_DAYS", "")
	t.Setenv("VAULTPULL_ROTATE_WARN_AGE_DAYS", "")
	t.Setenv("VAULTPULL_ROTATE_DRY_RUN", "")

	p := FromEnv()
	if p.MaxAge != 90*24*time.Hour {
		t.Errorf("MaxAge: got %v, want %v", p.MaxAge, 90*24*time.Hour)
	}
	if p.WarnAge != 75*24*time.Hour {
		t.Errorf("WarnAge: got %v, want %v", p.WarnAge, 75*24*time.Hour)
	}
	if p.DryRun {
		t.Error("DryRun should default to false")
	}
}

func TestFromEnvReadsMaxAgeDays(t *testing.T) {
	t.Setenv("VAULTPULL_ROTATE_MAX_AGE_DAYS", "60")
	t.Setenv("VAULTPULL_ROTATE_WARN_AGE_DAYS", "45")

	p := FromEnv()
	if p.MaxAge != 60*24*time.Hour {
		t.Errorf("MaxAge: got %v, want %v", p.MaxAge, 60*24*time.Hour)
	}
	if p.WarnAge != 45*24*time.Hour {
		t.Errorf("WarnAge: got %v, want %v", p.WarnAge, 45*24*time.Hour)
	}
}

func TestFromEnvReadsDryRun(t *testing.T) {
	t.Setenv("VAULTPULL_ROTATE_DRY_RUN", "true")

	p := FromEnv()
	if !p.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestFromEnvIgnoresInvalidMaxAge(t *testing.T) {
	t.Setenv("VAULTPULL_ROTATE_MAX_AGE_DAYS", "not-a-number")

	p := FromEnv()
	if p.MaxAge != 90*24*time.Hour {
		t.Errorf("MaxAge should fall back to default, got %v", p.MaxAge)
	}
}

func TestFromEnvClampsWarnAgeWhenAboveMaxAge(t *testing.T) {
	t.Setenv("VAULTPULL_ROTATE_MAX_AGE_DAYS", "30")
	t.Setenv("VAULTPULL_ROTATE_WARN_AGE_DAYS", "50")

	p := FromEnv()
	if p.WarnAge >= p.MaxAge {
		t.Errorf("WarnAge (%v) should be less than MaxAge (%v)", p.WarnAge, p.MaxAge)
	}
}
