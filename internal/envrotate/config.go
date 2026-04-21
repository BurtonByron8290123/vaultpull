package envrotate

import (
	"os"
	"strconv"
	"time"
)

// FromEnv builds a Policy from environment variables, falling back to defaults.
//
//	VAULTPULL_ROTATE_MAX_AGE_DAYS   – max age in days before rotation required (default 90)
//	VAULTPULL_ROTATE_WARN_AGE_DAYS  – warn age in days (default 75)
//	VAULTPULL_ROTATE_DRY_RUN        – if "true", report only without error (default false)
func FromEnv() Policy {
	p := DefaultPolicy()

	if v := os.Getenv("VAULTPULL_ROTATE_MAX_AGE_DAYS"); v != "" {
		if days, err := strconv.Atoi(v); err == nil && days > 0 {
			p.MaxAge = time.Duration(days) * 24 * time.Hour
		}
	}

	if v := os.Getenv("VAULTPULL_ROTATE_WARN_AGE_DAYS"); v != "" {
		if days, err := strconv.Atoi(v); err == nil && days > 0 {
			p.WarnAge = time.Duration(days) * 24 * time.Hour
		}
	}

	if os.Getenv("VAULTPULL_ROTATE_DRY_RUN") == "true" {
		p.DryRun = true
	}

	// Guard: if WarnAge was set independently and is now >= MaxAge, clamp it.
	if p.WarnAge >= p.MaxAge {
		p.WarnAge = p.MaxAge - 24*time.Hour
	}

	return p
}
