// Package ratelimit provides a simple token-bucket rate limiter for Vault API calls.
package ratelimit

import (
	"context"
	"fmt"
	"time"
)

// Policy defines rate limiting behaviour.
type Policy struct {
	// RequestsPerSecond is the maximum number of requests allowed per second.
	RequestsPerSecond float64
	// Burst is the maximum number of requests that can be made at once.
	Burst int
}

// DefaultPolicy returns a sensible default rate limit policy.
func DefaultPolicy() Policy {
	return Policy{
		RequestsPerSecond: 10,
		Burst:             5,
	}
}

// Validate returns an error if the policy is invalid.
func (p Policy) Validate() error {
	if p.RequestsPerSecond <= 0 {
		return fmt.Errorf("ratelimit: RequestsPerSecond must be > 0, got %f", p.RequestsPerSecond)
	}
	if p.Burst <= 0 {
		return fmt.Errorf("ratelimit: Burst must be > 0, got %d", p.Burst)
	}
	return nil
}

// Limiter controls the rate of outgoing requests.
type Limiter struct {
	tokens   float64
	max      float64
	rate     float64 // tokens per nanosecond
	last     time.Time
	waitFunc func(context.Context, time.Duration) error
}

// New creates a Limiter from the given Policy.
func New(p Policy) (*Limiter, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return &Limiter{
		tokens:   float64(p.Burst),
		max:      float64(p.Burst),
		rate:     p.RequestsPerSecond / float64(time.Second),
		last:     time.Now(),
		waitFunc: sleepContext,
	}, nil
}

// Wait blocks until a token is available or ctx is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	now := time.Now()
	elapsed := now.Sub(l.last).Nanoseconds()
	l.tokens += float64(elapsed) * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}
	l.last = now

	if l.tokens >= 1 {
		l.tokens--
		return nil
	}

	wait := time.Duration((1 - l.tokens) / l.rate)
	l.tokens = 0
	return l.waitFunc(ctx, wait)
}

func sleepContext(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
