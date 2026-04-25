// Package validate checks that fetched secret maps satisfy
// user-defined constraints before they are written to disk.
package validate

import (
	"errors"
	"fmt"
	"strings"
)

// Rule describes a single validation constraint.
type Rule struct {
	// Key is the exact secret key that must be present.
	Key string
	// MinLen is the minimum required value length (0 = no check).
	MinLen int
	// Disallow is a list of substrings that must not appear in the value.
	Disallow []string
}

// Policy holds all rules for a pull operation.
type Policy struct {
	Rules []Rule
}

// DefaultPolicy returns an empty policy that passes everything.
func DefaultPolicy() Policy { return Policy{} }

// Validator applies a Policy against a secret map.
type Validator struct {
	policy Policy
}

// New creates a Validator for the given policy.
func New(p Policy) *Validator { return &Validator{policy: p} }

// Validate checks secrets against every rule and returns a combined error
// listing all violations, or nil when all rules pass.
func (v *Validator) Validate(secrets map[string]string) error {
	var errs []string
	for _, r := range v.policy.Rules {
		val, ok := secrets[r.Key]
		if !ok {
			errs = append(errs, fmt.Sprintf("required key %q is missing", r.Key))
			continue
		}
		if r.MinLen > 0 && len(val) < r.MinLen {
			errs = append(errs, fmt.Sprintf("key %q: value length %d is below minimum %d", r.Key, len(val), r.MinLen))
		}
		for _, sub := range r.Disallow {
			if strings.Contains(val, sub) {
				errs = append(errs, fmt.Sprintf("key %q: value contains disallowed substring %q", r.Key, sub))
			}
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// AddRule appends a new Rule to the Policy and returns the Validator
// to allow method chaining.
func (v *Validator) AddRule(r Rule) *Validator {
	v.policy.Rules = append(v.policy.Rules, r)
	return v
}
