// Package envhealth checks the overall health of a local .env file
// by running a configurable set of assertions and reporting violations.
package envhealth

import (
	"fmt"
	"strings"
)

// Status represents the outcome of a health check.
type Status int

const (
	StatusOK      Status = iota
	StatusWarning        // non-fatal issues detected
	StatusError          // one or more required checks failed
)

// Violation describes a single failed assertion.
type Violation struct {
	Key     string
	Message string
}

// Report holds the aggregated result of all checks.
type Report struct {
	Status     Status
	Violations []Violation
}

// Summary returns a human-readable one-line summary.
func (r Report) Summary() string {
	switch r.Status {
	case StatusOK:
		return "env healthy: no issues found"
	case StatusWarning:
		return fmt.Sprintf("env warning: %d issue(s) detected", len(r.Violations))
	default:
		return fmt.Sprintf("env unhealthy: %d violation(s) found", len(r.Violations))
	}
}

// Policy controls which assertions are enabled.
type Policy struct {
	RequiredKeys    []string // keys that must be present and non-empty
	ForbiddenKeys   []string // keys that must not appear
	NoEmptyValues   bool     // warn when any value is empty
}

// DefaultPolicy returns a Policy with no assertions enabled.
func DefaultPolicy() Policy { return Policy{} }

// Checker runs health assertions against an env map.
type Checker struct {
	policy Policy
}

// New creates a Checker with the given Policy.
func New(p Policy) *Checker {
	return &Checker{policy: p}
}

// Check evaluates all assertions and returns a Report.
func (c *Checker) Check(env map[string]string) Report {
	var violations []Violation

	for _, key := range c.policy.RequiredKeys {
		v, ok := env[key]
		if !ok || strings.TrimSpace(v) == "" {
			violations = append(violations, Violation{
				Key:     key,
				Message: "required key is missing or empty",
			})
		}
	}

	for _, key := range c.policy.ForbiddenKeys {
		if _, ok := env[key]; ok {
			violations = append(violations, Violation{
				Key:     key,
				Message: "forbidden key is present",
			})
		}
	}

	if c.policy.NoEmptyValues {
		for k, v := range env {
			if strings.TrimSpace(v) == "" {
				violations = append(violations, Violation{
					Key:     k,
					Message: "value is empty",
				})
			}
		}
	}

	status := StatusOK
	if len(violations) > 0 {
		status = StatusError
	}
	return Report{Status: status, Violations: violations}
}
