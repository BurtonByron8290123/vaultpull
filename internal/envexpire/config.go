package envexpire

import (
	"os"
	"strconv"
	"time"
)

const (
	envMaxAgeDays  = "VAULTPULL_EXPIRE_MAX_AGE_DAYS"
	envWarnAgeDays = "VAULTPULL_EXPIRE_WARN_AGE_DAYS"
)

// FromEnv builds a Policy from environment variables, falling back to
// DefaultPolicy() for any value that is absent or invalid.
func FromEnv() Policy {
	p := DefaultPolicy()
	if v := os.Getenv(envMaxAgeDays); v != "" {
		if days, err := strconv.Atoi(v); err == nil && days > 0 {
			p.MaxAge = time.Duration(days) * 24 * time.Hour
		}
	}
	if v := os.Getenv(envWarnAgeDays); v != "" {
		if days, err := strconv.Atoi(v); err == nil && days > 0 {
			warn := time.Duration(days) * 24 * time.Hour
			if warn <= p.MaxAge {
				p.WarnAge = warn
			}
		}
	}
	return p
}
