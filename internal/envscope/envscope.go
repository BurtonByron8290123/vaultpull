// Package envscope restricts which keys may be written to a target env file
// based on a configurable allow/deny list of key prefixes.
package envscope

import (
	"fmt"
	"strings"
)

// Policy defines the allow and deny prefix rules.
type Policy struct {
	Allow []string // if non-empty, only keys matching a prefix are kept
	Deny  []string // keys matching any prefix are always removed
}

// Scoper applies a Policy to a map of env vars.
type Scoper struct {
	policy Policy
}

// New returns a Scoper for the given Policy.
func New(p Policy) (*Scoper, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	return &Scoper{policy: p}, nil
}

// Apply returns a new map containing only the keys permitted by the policy.
func (s *Scoper) Apply(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		if s.allowed(k) {
			out[k] = v
		}
	}
	return out
}

func (s *Scoper) allowed(key string) bool {
	for _, d := range s.policy.Deny {
		if strings.HasPrefix(key, d) {
			return false
		}
	}
	if len(s.policy.Allow) == 0 {
		return true
	}
	for _, a := range s.policy.Allow {
		if strings.HasPrefix(key, a) {
			return true
		}
	}
	return false
}

func validate(p Policy) error {
	for _, a := range p.Allow {
		if strings.TrimSpace(a) == "" {
			return fmt.Errorf("envscope: allow list contains blank prefix")
		}
	}
	for _, d := range p.Deny {
		if strings.TrimSpace(d) == "" {
			return fmt.Errorf("envscope: deny list contains blank prefix")
		}
	}
	return nil
}
