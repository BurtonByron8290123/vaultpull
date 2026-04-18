// Package backoff provides configurable exponential back-off policies
// used when retrying transient Vault errors.
package backoff

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// Policy holds the parameters for exponential back-off with jitter.
type Policy struct {
	InitialInterval time.Duration
	Multiplier      float64
	MaxInterval     time.Duration
	Jitter          float64 // fraction in [0,1]
}

// DefaultPolicy returns a sensible default back-off policy.
func DefaultPolicy() Policy {
	return Policy{
		InitialInterval: 200 * time.Millisecond,
		Multiplier:      2.0,
		MaxInterval:     30 * time.Second,
		Jitter:          0.2,
	}
}

// Validate returns an error if the policy contains invalid values.
func (p Policy) Validate() error {
	if p.InitialInterval <= 0 {
		return ErrInvalidInitial
	}
	if p.Multiplier < 1 {
		return ErrInvalidMultiplier
	}
	if p.MaxInterval < p.InitialInterval {
		return ErrInvalidMaxInterval
	}
	if p.Jitter < 0 || p.Jitter > 1 {
		return ErrInvalidJitter
	}
	return nil
}

// Next returns the back-off duration for the given attempt (0-indexed).
func (p Policy) Next(attempt int) time.Duration {
	base := float64(p.InitialInterval) * math.Pow(p.Multiplier, float64(attempt))
	if base > float64(p.MaxInterval) {
		base = float64(p.MaxInterval)
	}
	jitter := (rand.Float64()*2 - 1) * p.Jitter * base
	d := time.Duration(base + jitter)
	if d < 0 {
		d = 0
	}
	return d
}

// Sleep waits for the back-off duration or until ctx is cancelled.
func (p Policy) Sleep(ctx context.Context, attempt int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(p.Next(attempt)):
		return nil
	}
}
