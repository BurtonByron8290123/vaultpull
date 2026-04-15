package template

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// varPattern matches ${VAR_NAME} or $VAR_NAME style references.
var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Renderer resolves template variables in path strings using a set of
// variables and, optionally, the process environment.
type Renderer struct {
	vars      map[string]string
	useEnv    bool
}

// New creates a Renderer. When useEnv is true, OS environment variables
// are consulted as a fallback after vars.
func New(vars map[string]string, useEnv bool) *Renderer {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &Renderer{vars: copy, useEnv: useEnv}
}

// Render replaces all variable references in s with their resolved values.
// Returns an error if any reference cannot be resolved.
func (r *Renderer) Render(s string) (string, error) {
	var renderErr error
	result := varPattern.ReplaceAllStringFunc(s, func(match string) string {
		if renderErr != nil {
			return match
		}
		name := extractName(match)
		if v, ok := r.vars[name]; ok {
			return v
		}
		if r.useEnv {
			if v, ok := os.LookupEnv(name); ok {
				return v
			}
		}
		renderErr = fmt.Errorf("template: unresolved variable %q in %q", name, s)
		return match
	})
	if renderErr != nil {
		return "", renderErr
	}
	return result, nil
}

// RenderAll renders every value in paths, returning the resolved slice or
// the first error encountered.
func (r *Renderer) RenderAll(paths []string) ([]string, error) {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		resolved, err := r.Render(p)
		if err != nil {
			return nil, err
		}
		out = append(out, resolved)
	}
	return out, nil
}

// extractName strips the ${ } or $ sigil from a matched token.
func extractName(match string) string {
	match = strings.TrimPrefix(match, "${") 
	match = strings.TrimSuffix(match, "}")
	match = strings.TrimPrefix(match, "$")
	return match
}
