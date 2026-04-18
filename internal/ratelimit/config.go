package ratelimit

import (
	"fmt"
	"os"
	"strconv"
)

const (
	envRPS   = "VAULTPULL_RATE_RPS"
	envBurst = "VAULTPULL_RATE_BURST"
)

// FromEnv reads rate limit configuration from environment variables.
// VAULTPULL_RATE_RPS  – requests per second (float)
// VAULTPULL_RATE_BURST – burst size (int)
// Missing variables fall back to DefaultPolicy values.
func FromEnv() (Policy, error) {
	p := DefaultPolicy()

	if v := os.Getenv(envRPS); v != "" {
		rps, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return p, fmt.Errorf("ratelimit: invalid %s %q: %w", envRPS, v, err)
		}
		p.RequestsPerSecond = rps
	}

	if v := os.Getenv(envBurst); v != "" {
		burst, err := strconv.Atoi(v)
		if err != nil {
			return p, fmt.Errorf("ratelimit: invalid %s %q: %w", envBurst, v, err)
		}
		p.Burst = burst
	}

	if err := p.Validate(); err != nil {
		return p, err
	}
	return p, nil
}
