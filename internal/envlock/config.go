package envlock

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultTimeout  = 10 * time.Second
	defaultLockExt  = ".lock"
	envLockTimeout  = "VAULTPULL_LOCK_TIMEOUT_SEC"
	envLockSuffix   = "VAULTPULL_LOCK_SUFFIX"
)

// Config controls envlock behaviour.
type Config struct {
	// Timeout is how long Acquire will wait before returning an error.
	Timeout time.Duration
	// Suffix is appended to the env file path to form the lock file path.
	Suffix string
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Timeout: defaultTimeout,
		Suffix:  defaultLockExt,
	}
}

// FromEnv builds a Config from environment variables, falling back to
// DefaultConfig for any value that is absent or invalid.
func FromEnv() Config {
	cfg := DefaultConfig()

	if raw := os.Getenv(envLockTimeout); raw != "" {
		if secs, err := strconv.Atoi(raw); err == nil && secs > 0 {
			cfg.Timeout = time.Duration(secs) * time.Second
		}
	}

	if suffix := os.Getenv(envLockSuffix); suffix != "" {
		cfg.Suffix = suffix
	}

	return cfg
}
