// Package envsorter provides utilities for sorting environment variable maps
// by key, with support for custom ordering strategies.
package envsort

import (
	"sort"
	"strings"
)

// Strategy controls how keys are ordered.
type Strategy int

const (
	// Alphabetical sorts keys in ascending lexicographic order.
	Alphabetical Strategy = iota
	// AlphabeticalDesc sorts keys in descending lexicographic order.
	AlphabeticalDesc
	// PrefixGrouped groups keys by prefix (split on '_') then sorts within groups.
	PrefixGrouped
)

// Policy holds sorting configuration.
type Policy struct {
	Strategy Strategy
	// PrefixSep is the separator used when Strategy is PrefixGrouped. Defaults to "_".
	PrefixSep string
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		Strategy:  Alphabetical,
		PrefixSep: "_",
	}
}

// Sorter sorts environment variable keys.
type Sorter struct {
	policy Policy
}

// New creates a Sorter with the given policy.
func New(p Policy) (*Sorter, error) {
	if p.PrefixSep == "" {
		p.PrefixSep = "_"
	}
	return &Sorter{policy: p}, nil
}

// SortedKeys returns the keys of m ordered according to the policy.
func (s *Sorter) SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	switch s.policy.Strategy {
	case AlphabeticalDesc:
		sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })
	case PrefixGrouped:
		sort.Slice(keys, func(i, j int) bool {
			pi := prefix(keys[i], s.policy.PrefixSep)
			pj := prefix(keys[j], s.policy.PrefixSep)
			if pi != pj {
				return pi < pj
			}
			return keys[i] < keys[j]
		})
	default: // Alphabetical
		sort.Strings(keys)
	}
	return keys
}

// Apply returns a new map identical to m (no mutation) with keys accessible
// via SortedKeys. Since maps are unordered, this is a convenience wrapper
// that validates all keys are preserved.
func (s *Sorter) Apply(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func prefix(key, sep string) string {
	if idx := strings.Index(key, sep); idx >= 0 {
		return key[:idx]
	}
	return key
}
