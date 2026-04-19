package leasetrack

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	envWarnThresholdSec = "VAULTPULL_LEASE_WARN_SEC"
)

// FromEnv builds a Policy from environment variables, falling back to
// DefaultPolicy() for any value that is absent or invalid.
func FromEnv() Policy {
	p := DefaultPolicy()

	if raw := os.Getenv(envWarnThresholdSec); raw != "" {
		if secs, err := strconv.Atoi(raw); err == nil && secs > 0 {
			p.WarnThreshold = time.Duration(secs) * time.Second
		}
	}

	return p
}

// FromEnvWithOverrides merges explicit overrides on top of the environment-
// derived policy. A zero-value duration in overrides is ignored.
func FromEnvWithOverrides(warn time.Duration) (Policy, error) {
	p := FromEnv()
	if warn > 0 {
		p.WarnThreshold = warn
	}
	if err := p.Validate(); err != nil {
		return Policy{}, fmt.Errorf("leasetrack config: %w", err)
	}
	return p, nil
}
