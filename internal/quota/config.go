package quota

import (
	"fmt"
	"os"
	"strconv"
)

const (
	envMaxRequests = "VAULTPULL_QUOTA_MAX_REQUESTS"
	defaultMax     = 10
)

// FromEnv builds a Policy from environment variables, falling back to
// defaults for missing or invalid values.
func FromEnv() (Policy, error) {
	p := DefaultPolicy()
	if raw := os.Getenv(envMaxRequests); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return Policy{}, fmt.Errorf("quota: invalid %s %q: %w", envMaxRequests, raw, err)
		}
		p.MaxRequestsPerPath = v
	}
	if err := p.Validate(); err != nil {
		return Policy{}, err
	}
	return p, nil
}
