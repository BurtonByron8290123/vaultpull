// Package envrotate provides secret rotation detection and enforcement
// for .env files synced from Vault.
package envrotate

import (
	"errors"
	"fmt"
	"time"
)

// ErrRotationRequired is returned when a secret has exceeded its max age.
var ErrRotationRequired = errors.New("envrotate: secret rotation required")

// Policy controls rotation behaviour.
type Policy struct {
	// MaxAge is the maximum allowed age of a secret before rotation is required.
	MaxAge time.Duration
	// WarnAge triggers a warning when a secret is older than this threshold.
	WarnAge time.Duration
	// DryRun reports rotation needs without returning an error.
	DryRun bool
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		MaxAge:  90 * 24 * time.Hour,
		WarnAge: 75 * 24 * time.Hour,
	}
}

// Rotator checks whether secrets require rotation based on their age.
type Rotator struct {
	policy Policy
	clock  func() time.Time
}

// New creates a Rotator with the given policy.
func New(p Policy) (*Rotator, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	return &Rotator{policy: p, clock: time.Now}, nil
}

// Status describes the rotation status of a secret.
type Status int

const (
	StatusOK      Status = iota
	StatusWarning        // approaching max age
	StatusExpired        // exceeded max age
)

// Result holds the outcome of a rotation check for a single key.
type Result struct {
	Key       string
	Age       time.Duration
	Status    Status
	RotatedAt time.Time
}

// Check evaluates a map of key → last-rotated timestamps and returns results.
func (r *Rotator) Check(rotatedAt map[string]time.Time) ([]Result, error) {
	now := r.clock()
	results := make([]Result, 0, len(rotatedAt))
	var expired []string

	for key, ts := range rotatedAt {
		age := now.Sub(ts)
		res := Result{Key: key, Age: age, RotatedAt: ts}
		switch {
		case age >= r.policy.MaxAge:
			res.Status = StatusExpired
			expired = append(expired, key)
		case age >= r.policy.WarnAge:
			res.Status = StatusWarning
		default:
			res.Status = StatusOK
		}
		results = append(results, res)
	}

	if len(expired) > 0 && !r.policy.DryRun {
		return results, fmt.Errorf("%w: keys=%v", ErrRotationRequired, expired)
	}
	return results, nil
}

func validate(p Policy) error {
	if p.MaxAge <= 0 {
		return errors.New("envrotate: MaxAge must be positive")
	}
	if p.WarnAge <= 0 {
		return errors.New("envrotate: WarnAge must be positive")
	}
	if p.WarnAge >= p.MaxAge {
		return errors.New("envrotate: WarnAge must be less than MaxAge")
	}
	return nil
}
