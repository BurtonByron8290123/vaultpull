// Package jitter adds randomised jitter to durations, useful for
// spreading out retry and polling intervals across many clients.
package jitter

import (
	"fmt"
	"math/rand"
	"time"
)

// Policy controls how jitter is applied to a base duration.
type Policy struct {
	// Factor is the maximum fraction of the base duration to add as jitter.
	// Must be in the range (0, 1].
	Factor float64
}

// DefaultPolicy returns a Policy with a 20 % jitter factor.
func DefaultPolicy() Policy {
	return Policy{Factor: 0.20}
}

// Validate returns an error if the policy is misconfigured.
func (p Policy) Validate() error {
	if p.Factor <= 0 || p.Factor > 1 {
		return fmt.Errorf("jitter: factor must be in (0, 1], got %v", p.Factor)
	}
	return nil
}

// Apply returns base + a random duration in [0, base*Factor).
func (p Policy) Apply(base time.Duration) time.Duration {
	max := float64(base) * p.Factor
	offset := time.Duration(rand.Float64() * max) //nolint:gosec
	return base + offset
}

// ApplyN calls Apply n times and returns all results.
func (p Policy) ApplyN(base time.Duration, n int) []time.Duration {
	out := make([]time.Duration, n)
	for i := range out {
		out[i] = p.Apply(base)
	}
	return out
}
