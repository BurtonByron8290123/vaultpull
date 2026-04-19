// Package scope restricts which Vault secret paths a pull operation may access.
package scope

import (
	"errors"
	"fmt"
	"strings"
)

// Policy defines allowed and denied path prefixes.
type Policy struct {
	Allow []string
	Deny  []string
}

// Scope enforces a Policy against secret paths.
type Scope struct {
	policy Policy
}

// New returns a Scope for the given Policy.
func New(p Policy) (*Scope, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	return &Scope{policy: p}, nil
}

// Allow reports whether path is permitted under the policy.
func (s *Scope) Allow(path string) bool {
	for _, d := range s.policy.Deny {
		if strings.HasPrefix(path, d) {
			return false
		}
	}
	if len(s.policy.Allow) == 0 {
		return true
	}
	for _, a := range s.policy.Allow {
		if strings.HasPrefix(path, a) {
			return true
		}
	}
	return false
}

// Filter returns only the paths permitted by the policy.
func (s *Scope) Filter(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if s.Allow(p) {
			out = append(out, p)
		}
	}
	return out
}

func validate(p Policy) error {
	for _, a := range p.Allow {
		if strings.TrimSpace(a) == "" {
			return errors.New("scope: allow entry must not be blank")
		}
	}
	for _, d := range p.Deny {
		if strings.TrimSpace(d) == "" {
			return errors.New("scope: deny entry must not be blank")
		}
	}
	for _, a := range p.Allow {
		for _, d := range p.Deny {
			if a == d {
				return fmt.Errorf("scope: path %q appears in both allow and deny", a)
			}
		}
	}
	return nil
}
