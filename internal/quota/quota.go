// Package quota enforces per-path request limits during a pull run.
package quota

import (
	"errors"
	"fmt"
	"sync"
)

// ErrQuotaExceeded is returned when the request limit for a path is reached.
var ErrQuotaExceeded = errors.New("quota exceeded")

// Policy defines the quota configuration.
type Policy struct {
	// MaxRequestsPerPath is the maximum number of Vault fetches allowed per path.
	MaxRequestsPerPath int
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{MaxRequestsPerPath: 10}
}

// Validate returns an error if the policy is invalid.
func (p Policy) Validate() error {
	if p.MaxRequestsPerPath <= 0 {
		return fmt.Errorf("quota: MaxRequestsPerPath must be > 0, got %d", p.MaxRequestsPerPath)
	}
	return nil
}

// Tracker counts requests per path and enforces the policy.
type Tracker struct {
	mu      sync.Mutex
	counts  map[string]int
	policy  Policy
}

// New creates a new Tracker with the given policy.
func New(p Policy) (*Tracker, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return &Tracker{counts: make(map[string]int), policy: p}, nil
}

// Allow records a request for path and returns ErrQuotaExceeded if the limit
// has been reached.
func (t *Tracker) Allow(path string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts[path]++
	if t.counts[path] > t.policy.MaxRequestsPerPath {
		return fmt.Errorf("%w: path %q reached limit of %d", ErrQuotaExceeded, path, t.policy.MaxRequestsPerPath)
	}
	return nil
}

// Count returns the current request count for path.
func (t *Tracker) Count(path string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.counts[path]
}

// Reset clears all counters.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts = make(map[string]int)
}
