// Package transform applies key/value transformations to secret maps
// before they are written to .env files.
package transform

import (
	"strings"
)

// Policy controls how keys and values are transformed.
type Policy struct {
	// PrefixStrip removes this prefix from every key, if present.
	PrefixStrip string
	// KeyCase converts keys: "upper", "lower", or "" (no change).
	KeyCase string
	// ValueTrimSpace trims leading/trailing whitespace from values.
	ValueTrimSpace bool
}

// Transformer applies a Policy to secret maps.
type Transformer struct {
	p Policy
}

// New returns a Transformer for the given Policy.
func New(p Policy) (*Transformer, error) {
	if p.KeyCase != "" && p.KeyCase != "upper" && p.KeyCase != "lower" {
		return nil, &InvalidKeyCaseError{Value: p.KeyCase}
	}
	return &Transformer{p: p}, nil
}

// Apply returns a new map with all transformations applied.
func (t *Transformer) Apply(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		k = t.transformKey(k)
		v = t.transformValue(v)
		out[k] = v
	}
	return out
}

func (t *Transformer) transformKey(k string) string {
	if t.p.PrefixStrip != "" {
		k = strings.TrimPrefix(k, t.p.PrefixStrip)
	}
	switch t.p.KeyCase {
	case "upper":
		k = strings.ToUpper(k)
	case "lower":
		k = strings.ToLower(k)
	}
	return k
}

func (t *Transformer) transformValue(v string) string {
	if t.p.ValueTrimSpace {
		v = strings.TrimSpace(v)
	}
	return v
}

// InvalidKeyCaseError is returned when an unknown KeyCase is specified.
type InvalidKeyCaseError struct {
	Value string
}

func (e *InvalidKeyCaseError) Error() string {
	return "transform: invalid key_case \"" + e.Value + "\": must be \"upper\" or \"lower\""
}
