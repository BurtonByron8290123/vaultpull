// Package envcheck verifies that required environment variables are present
// and optionally non-empty before a pull operation proceeds.
package envcheck

import (
	"errors"
	"fmt"
	"strings"
)

// Policy holds the configuration for the environment checker.
type Policy struct {
	// Required lists keys that must be present in the env map.
	Required []string
	// NonEmpty lists keys that must be present and have a non-empty value.
	NonEmpty []string
}

// Violation describes a single failed check.
type Violation struct {
	Key    string
	Reason string
}

func (v Violation) Error() string {
	return fmt.Sprintf("envcheck: key %q %s", v.Key, v.Reason)
}

// Checker validates an env map against a Policy.
type Checker struct {
	policy Policy
}

// New returns a Checker for the given policy.
func New(p Policy) (*Checker, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	return &Checker{policy: p}, nil
}

// Check inspects env and returns all violations found.
// A nil error means all checks passed.
func (c *Checker) Check(env map[string]string) error {
	var errs []string

	for _, key := range c.policy.Required {
		if _, ok := env[key]; !ok {
			errs = append(errs, Violation{Key: key, Reason: "is required but missing"}.Error())
		}
	}

	for _, key := range c.policy.NonEmpty {
		v, ok := env[key]
		if !ok {
			errs = append(errs, Violation{Key: key, Reason: "is required but missing"}.Error())
		} else if strings.TrimSpace(v) == "" {
			errs = append(errs, Violation{Key: key, Reason: "must not be empty"}.Error())
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errs, "; "))
}

func validate(p Policy) error {
	for _, k := range p.Required {
		if strings.TrimSpace(k) == "" {
			return errors.New("envcheck: required key must not be blank")
		}
	}
	for _, k := range p.NonEmpty {
		if strings.TrimSpace(k) == "" {
			return errors.New("envcheck: non-empty key must not be blank")
		}
	}
	return nil
}
