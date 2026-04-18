package expire

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	envTTLSeconds = "VAULTPULL_SECRET_TTL_SECONDS"
	defaultTTLSec = 300
)

// FromEnv builds a Policy from environment variables, falling back to defaults.
func FromEnv() (Policy, error) {
	p := DefaultPolicy()

	if raw := os.Getenv(envTTLSeconds); raw != "" {
		secs, err := strconv.Atoi(raw)
		if err != nil {
			return Policy{}, fmt.Errorf("expire: invalid %s value %q: %w", envTTLSeconds, raw, err)
		}
		if secs <= 0 {
			return Policy{}, fmt.Errorf("expire: %s must be positive, got %d", envTTLSeconds, secs)
		}
		p.TTL = time.Duration(secs) * time.Second
	}

	return p, nil
}
