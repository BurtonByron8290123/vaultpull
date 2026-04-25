// Package envexpire checks whether secrets in a .env file have exceeded
// their maximum allowed age and reports which keys need to be refreshed.
package envexpire

import (
	"errors"
	"fmt"
	"time"
)

// Status represents the expiry state of a single key.
type Status int

const (
	StatusOK      Status = iota // within max age
	StatusWarning               // approaching max age
	StatusExpired               // exceeded max age
)

// Result holds the expiry result for one key.
type Result struct {
	Key    string
	Age    time.Duration
	Status Status
}

// Policy controls expiry thresholds.
type Policy struct {
	MaxAge  time.Duration // keys older than this are expired
	WarnAge time.Duration // keys older than this trigger a warning
}

// DefaultPolicy returns a sensible default expiry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAge:  30 * 24 * time.Hour,
		WarnAge: 25 * 24 * time.Hour,
	}
}

func (p Policy) validate() error {
	if p.MaxAge <= 0 {
		return errors.New("envexpire: MaxAge must be positive")
	}
	if p.WarnAge <= 0 {
		return errors.New("envexpire: WarnAge must be positive")
	}
	if p.WarnAge > p.MaxAge {
		return fmt.Errorf("envexpire: WarnAge (%v) must not exceed MaxAge (%v)", p.WarnAge, p.MaxAge)
	}
	return nil
}

// Checker evaluates key ages against a Policy.
type Checker struct {
	policy Policy
	clock  func() time.Time
}

// New creates a Checker with the given policy.
func New(p Policy) (*Checker, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &Checker{policy: p, clock: time.Now}, nil
}

// Check evaluates a map of key→timestamp (RFC3339) and returns per-key results.
func (c *Checker) Check(timestamps map[string]time.Time) []Result {
	now := c.clock()
	results := make([]Result, 0, len(timestamps))
	for key, ts := range timestamps {
		age := now.Sub(ts)
		status := StatusOK
		switch {
		case age >= c.policy.MaxAge:
			status = StatusExpired
		case age >= c.policy.WarnAge:
			status = StatusWarning
		}
		results = append(results, Result{Key: key, Age: age, Status: status})
	}
	return results
}

// Expired returns only the results whose status is StatusExpired.
func Expired(results []Result) []Result {
	out := results[:0:0]
	for _, r := range results {
		if r.Status == StatusExpired {
			out = append(out, r)
		}
	}
	return out
}
