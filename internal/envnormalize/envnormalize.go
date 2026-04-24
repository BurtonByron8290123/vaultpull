// Package envnormalize provides key and value normalization for env maps.
// It ensures consistent casing, trimming, and formatting before secrets
// are written to .env files.
package envnormalize

import (
	"strings"
)

// Strategy controls how keys are normalized.
type Strategy int

const (
	// StrategyUpper converts all keys to UPPER_CASE.
	StrategyUpper Strategy = iota
	// StrategyLower converts all keys to lower_case.
	StrategyLower
	// StrategyNone leaves keys unchanged.
	StrategyNone
)

// Policy defines normalization behaviour.
type Policy struct {
	KeyStrategy   Strategy
	TrimValues    bool
	CollapseEmpty bool // replace empty values with a zero-length string (no-op placeholder)
}

// DefaultPolicy returns a sensible default: uppercase keys, trimmed values.
func DefaultPolicy() Policy {
	return Policy{
		KeyStrategy: StrategyUpper,
		TrimValues:  true,
	}
}

// Normalizer applies a Policy to an env map.
type Normalizer struct {
	policy Policy
}

// New creates a Normalizer with the supplied policy.
func New(p Policy) (*Normalizer, error) {
	return &Normalizer{policy: p}, nil
}

// Apply returns a new map with all keys and values normalized according
// to the policy. The original map is never mutated.
func (n *Normalizer) Apply(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		nk := n.normalizeKey(k)
		nv := n.normalizeValue(v)
		out[nk] = nv
	}
	return out
}

func (n *Normalizer) normalizeKey(k string) string {
	switch n.policy.KeyStrategy {
	case StrategyUpper:
		return strings.ToUpper(k)
	case StrategyLower:
		return strings.ToLower(k)
	default:
		return k
	}
}

func (n *Normalizer) normalizeValue(v string) string {
	if n.policy.TrimValues {
		v = strings.TrimSpace(v)
	}
	return v
}
