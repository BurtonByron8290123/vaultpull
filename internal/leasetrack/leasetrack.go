// Package leasetrack monitors Vault secret lease expiry and warns when
// secrets are approaching their TTL.
package leasetrack

import (
	"fmt"
	"time"
)

// Entry holds lease metadata for a single secret path.
type Entry struct {
	Path      string
	LeaseID   string
	IssuedAt  time.Time
	LeaseTTL  time.Duration
}

// ExpiresAt returns the absolute expiry time for the entry.
func (e Entry) ExpiresAt() time.Time {
	return e.IssuedAt.Add(e.LeaseTTL)
}

// RemainingTTL returns how long until the lease expires relative to now.
func (e Entry) RemainingTTL(now time.Time) time.Duration {
	return e.ExpiresAt().Sub(now)
}

// Status describes the freshness of a lease.
type Status int

const (
	StatusFresh   Status = iota
	StatusWarning        // within warn threshold
	StatusExpired
)

// Policy controls thresholds for lease status evaluation.
type Policy struct {
	WarnThreshold time.Duration // warn when remaining TTL falls below this
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{WarnThreshold: 10 * time.Minute}
}

// Validate returns an error if the policy is misconfigured.
func (p Policy) Validate() error {
	if p.WarnThreshold <= 0 {
		return fmt.Errorf("leasetrack: WarnThreshold must be positive")
	}
	return nil
}

// Tracker evaluates lease entries against a policy.
type Tracker struct {
	policy Policy
	clock  func() time.Time
}

// New creates a Tracker using the supplied policy.
func New(p Policy) (*Tracker, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return &Tracker{policy: p, clock: time.Now}, nil
}

// Check returns the Status for a given Entry.
func (t *Tracker) Check(e Entry) Status {
	now := t.clock()
	if now.After(e.ExpiresAt()) {
		return StatusExpired
	}
	if e.RemainingTTL(now) <= t.policy.WarnThreshold {
		return StatusWarning
	}
	return StatusFresh
}

// CheckAll evaluates a slice of entries and returns only those that are
// warning or expired.
func (t *Tracker) CheckAll(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if t.Check(e) != StatusFresh {
			out = append(out, e)
		}
	}
	return out
}
