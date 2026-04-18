// Package expire provides TTL-based expiry checking for cached secrets.
package expire

import (
	"time"
)

// Policy defines expiry behaviour for a secret entry.
type Policy struct {
	TTL        time.Duration
	ClockFunc  func() time.Time
}

// DefaultPolicy returns a Policy with a 5-minute TTL.
func DefaultPolicy() Policy {
	return Policy{
		TTL:       5 * time.Minute,
		ClockFunc: time.Now,
	}
}

// Validate returns an error if the policy is misconfigured.
func (p Policy) Validate() error {
	if p.TTL <= 0 {
		return ErrInvalidTTL
	}
	if p.ClockFunc == nil {
		return ErrNilClock
	}
	return nil
}

// IsExpired reports whether fetchedAt is older than the policy TTL.
func (p Policy) IsExpired(fetchedAt time.Time) bool {
	return p.ClockFunc().After(fetchedAt.Add(p.TTL))
}

// ExpiresAt returns the absolute expiry time for a given fetch time.
func (p Policy) ExpiresAt(fetchedAt time.Time) time.Time {
	return fetchedAt.Add(p.TTL)
}
