// Package envfreeze provides a mechanism to freeze (write-protect) a set of
// environment variable keys so they cannot be overwritten during a pull.
package envfreeze

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// ErrFrozenKey is returned when an attempt is made to overwrite a frozen key.
var ErrFrozenKey = errors.New("envfreeze: key is frozen")

// Policy holds the configuration for the freeze behaviour.
type Policy struct {
	// Keys is the explicit set of keys that must not be overwritten.
	Keys []string
	// DryRun reports violations without returning an error.
	DryRun bool
}

// Freezer enforces the frozen-key policy against an incoming map.
type Freezer struct {
	policy Policy
	frozen map[string]struct{}
}

// DefaultPolicy returns a Policy with no frozen keys and DryRun disabled.
func DefaultPolicy() Policy {
	return Policy{}
}

// New creates a Freezer from the supplied Policy.
func New(p Policy) (*Freezer, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	frozen := make(map[string]struct{}, len(p.Keys))
	for _, k := range p.Keys {
		frozen[strings.ToUpper(strings.TrimSpace(k))] = struct{}{}
	}
	return &Freezer{policy: p, frozen: frozen}, nil
}

// Check inspects incoming against current and returns an error (or prints a
// warning when DryRun is true) for every frozen key that would change value.
func (f *Freezer) Check(current, incoming map[string]string) error {
	var violations []string
	for k := range f.frozen {
		newVal, exists := incoming[k]
		if !exists {
			continue
		}
		oldVal, had := current[k]
		if had && oldVal != newVal {
			violations = append(violations, k)
		}
	}
	if len(violations) == 0 {
		return nil
	}
	msg := fmt.Sprintf("%w: %s", ErrFrozenKey, strings.Join(violations, ", "))
	if f.policy.DryRun {
		fmt.Fprintf(os.Stderr, "[envfreeze dry-run] %s\n", msg)
		return nil
	}
	return errors.New(msg)
}

// IsFrozen reports whether key k is in the frozen set.
func (f *Freezer) IsFrozen(k string) bool {
	_, ok := f.frozen[strings.ToUpper(strings.TrimSpace(k))]
	return ok
}

func validate(p Policy) error {
	for _, k := range p.Keys {
		if strings.TrimSpace(k) == "" {
			return errors.New("envfreeze: frozen key must not be blank")
		}
	}
	return nil
}
