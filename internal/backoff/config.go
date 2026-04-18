package backoff

import (
	"os"
	"strconv"
	"time"
)

// FromEnv builds a Policy from environment variables, falling back to
// DefaultPolicy() for any value that is absent or unparseable.
//
//	VAULTPULL_BACKOFF_INITIAL_MS  – initial interval in milliseconds
//	VAULTPULL_BACKOFF_MULTIPLIER  – growth multiplier (float)
//	VAULTPULL_BACKOFF_MAX_MS      – maximum interval in milliseconds
//	VAULTPULL_BACKOFF_JITTER      – jitter fraction in [0,1]
func FromEnv() Policy {
	p := DefaultPolicy()
	if v := os.Getenv("VAULTPULL_BACKOFF_INITIAL_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil && ms > 0 {
			p.InitialInterval = time.Duration(ms) * time.Millisecond
		}
	}
	if v := os.Getenv("VAULTPULL_BACKOFF_MULTIPLIER"); v != "" {
		if m, err := strconv.ParseFloat(v, 64); err == nil && m >= 1 {
			p.Multiplier = m
		}
	}
	if v := os.Getenv("VAULTPULL_BACKOFF_MAX_MS"); v != "" {
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil && ms > 0 {
			p.MaxInterval = time.Duration(ms) * time.Millisecond
		}
	}
	if v := os.Getenv("VAULTPULL_BACKOFF_JITTER"); v != "" {
		if j, err := strconv.ParseFloat(v, 64); err == nil && j >= 0 && j <= 1 {
			p.Jitter = j
		}
	}
	return p
}
