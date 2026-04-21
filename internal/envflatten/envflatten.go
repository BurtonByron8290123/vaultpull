// Package envflatten provides utilities for flattening nested secret maps
// (e.g. from Vault KV responses) into a flat KEY=VALUE env map.
package envflatten

import (
	"fmt"
	"strings"
)

// Policy controls how nested keys are flattened.
type Policy struct {
	// Separator is inserted between path segments. Defaults to "_".
	Separator string
	// UpperCase converts all keys to uppercase when true.
	UpperCase bool
	// Prefix is prepended to every flattened key.
	Prefix string
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		Separator: "_",
		UpperCase: true,
	}
}

// Flattener flattens nested maps into a single-level env map.
type Flattener struct {
	policy Policy
}

// New creates a Flattener with the given policy.
func New(p Policy) (*Flattener, error) {
	if p.Separator == "" {
		return nil, fmt.Errorf("envflatten: separator must not be empty")
	}
	return &Flattener{policy: p}, nil
}

// Flatten converts a nested map[string]any into a flat map[string]string.
// Nested maps are recursed; all other values are converted via fmt.Sprintf.
func (f *Flattener) Flatten(input map[string]any) map[string]string {
	out := make(map[string]string)
	f.flatten(input, f.policy.Prefix, out)
	return out
}

func (f *Flattener) flatten(input map[string]any, prefix string, out map[string]string) {
	for k, v := range input {
		key := f.buildKey(prefix, k)
		switch child := v.(type) {
		case map[string]any:
			f.flatten(child, key, out)
		case map[string]string:
			for ck, cv := range child {
				out[f.buildKey(key, ck)] = cv
			}
		default:
			out[key] = fmt.Sprintf("%v", v)
		}
	}
}

func (f *Flattener) buildKey(prefix, segment string) string {
	var key string
	if prefix == "" {
		key = segment
	} else {
		key = prefix + f.policy.Separator + segment
	}
	if f.policy.UpperCase {
		return strings.ToUpper(key)
	}
	return key
}
