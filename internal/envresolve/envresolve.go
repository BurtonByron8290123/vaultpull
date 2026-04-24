// Package envresolve expands variable references within env map values,
// substituting ${VAR} or $VAR placeholders with values from the same map
// or from a supplied override set.
package envresolve

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var placeholderRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Policy controls resolver behaviour.
type Policy struct {
	// AllowEnvFallback allows falling back to the process environment when a
	// key is absent from both the input map and overrides.
	AllowEnvFallback bool
	// ErrorOnMissing returns an error when a placeholder cannot be resolved.
	ErrorOnMissing bool
}

// DefaultPolicy returns a sensible default policy.
func DefaultPolicy() Policy {
	return Policy{
		AllowEnvFallback: true,
		ErrorOnMissing:   false,
	}
}

// Resolver expands variable references inside env values.
type Resolver struct {
	policy    Policy
	overrides map[string]string
}

// New creates a Resolver with the given policy and optional overrides.
func New(p Policy, overrides map[string]string) *Resolver {
	ov := make(map[string]string, len(overrides))
	for k, v := range overrides {
		ov[k] = v
	}
	return &Resolver{policy: p, overrides: ov}
}

// Apply returns a new map with all placeholder references expanded.
func (r *Resolver) Apply(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		expanded, err := r.expand(v, env)
		if err != nil {
			return nil, fmt.Errorf("envresolve: key %q: %w", k, err)
		}
		out[k] = expanded
	}
	return out, nil
}

func (r *Resolver) expand(value string, env map[string]string) (string, error) {
	var resolveErr error
	result := placeholderRe.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		name := extractName(match)
		if v, ok := r.overrides[name]; ok {
			return v
		}
		if v, ok := env[name]; ok {
			return v
		}
		if r.policy.AllowEnvFallback {
			if v, ok := os.LookupEnv(name); ok {
				return v
			}
		}
		if r.policy.ErrorOnMissing {
			resolveErr = fmt.Errorf("unresolved placeholder %q", name)
			return match
		}
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

func extractName(placeholder string) string {
	s := strings.TrimPrefix(placeholder, "$")
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	return s
}
