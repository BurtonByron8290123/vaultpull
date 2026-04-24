// Package envtemplate renders environment variable maps using Go text/template
// syntax, allowing dynamic construction of values from other keys or overrides.
package envtemplate

import (
	"bytes"
	"fmt"
	"text/template"
)

// Policy controls how template rendering behaves.
type Policy struct {
	// ErrorOnMissing causes rendering to fail when a referenced key is absent.
	ErrorOnMissing bool
	// LeftDelim and RightDelim override the default {{ }} delimiters.
	LeftDelim  string
	RightDelim string
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		ErrorOnMissing: false,
		LeftDelim:      "{{",
		RightDelim:     "}}",
	}
}

// Renderer applies Go template rendering to env map values.
type Renderer struct {
	policy Policy
}

// New returns a Renderer configured with the given Policy.
func New(p Policy) (*Renderer, error) {
	if p.LeftDelim == "" || p.RightDelim == "" {
		return nil, fmt.Errorf("envtemplate: delimiters must not be empty")
	}
	return &Renderer{policy: p}, nil
}

// Apply renders each value in src as a template, with src itself as the data
// context. Returns a new map; src is never mutated.
func (r *Renderer) Apply(src map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		rendered, err := r.render(v, src)
		if err != nil {
			return nil, fmt.Errorf("envtemplate: key %q: %w", k, err)
		}
		out[k] = rendered
	}
	return out, nil
}

func (r *Renderer) render(text string, data map[string]string) (string, error) {
	option := "missingkey=zero"
	if r.policy.ErrorOnMissing {
		option = "missingkey=error"
	}
	tmpl, err := template.New("").
		Delims(r.policy.LeftDelim, r.policy.RightDelim).
		Option(option).
		Parse(text)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
