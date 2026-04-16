// Package resolve handles resolving secret paths from config variables and environment overrides.
package resolve

import (
	"fmt"
	"os"
	"strings"
)

// Resolver expands path templates using environment variables and explicit vars.
type Resolver struct {
	vars map[string]string
}

// New returns a Resolver seeded with the given vars. Environment variables
// are consulted as a fallback when a placeholder is not found in vars.
func New(vars map[string]string) *Resolver {
	merged := make(map[string]string, len(vars))
	for k, v := range vars {
		merged[k] = v
	}
	return &Resolver{vars: merged}
}

// Resolve expands all ${VAR} and $VAR placeholders in path.
// Returns an error if any placeholder cannot be resolved.
func (r *Resolver) Resolve(path string) (string, error) {
	var missing []string
	result := os.Expand(path, func(key string) string {
		if v, ok := r.vars[key]; ok {
			return v
		}
		if v, ok := os.LookupEnv(key); ok {
			return v
		}
		missing = append(missing, key)
		return ""
	})
	if len(missing) > 0 {
		return "", fmt.Errorf("resolve: unresolved placeholders: %s", strings.Join(missing, ", "))
	}
	return result, nil
}

// ResolveAll expands each path in the slice, returning the first error encountered.
func (r *Resolver) ResolveAll(paths []string) ([]string, error) {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		resolved, err := r.Resolve(p)
		if err != nil {
			return nil, err
		}
		out = append(out, resolved)
	}
	return out, nil
}
