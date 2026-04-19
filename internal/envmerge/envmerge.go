// Package envmerge provides strategies for merging multiple env maps
// from different Vault paths into a single resolved map.
package envmerge

import "fmt"

// Strategy controls how key conflicts are resolved.
type Strategy int

const (
	// StrategyLastWins uses the value from the last source that defines the key.
	StrategyLastWins Strategy = iota
	// StrategyFirstWins keeps the value from the first source that defines the key.
	StrategyFirstWins
	// StrategyError returns an error when two sources define the same key with different values.
	StrategyError
)

// Policy holds configuration for the merge operation.
type Policy struct {
	Strategy Strategy
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{Strategy: StrategyLastWins}
}

// Merger merges multiple env maps according to a policy.
type Merger struct {
	policy Policy
}

// New creates a Merger with the given policy.
func New(p Policy) (*Merger, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	return &Merger{policy: p}, nil
}

func validate(p Policy) error {
	if p.Strategy < StrategyLastWins || p.Strategy > StrategyError {
		return fmt.Errorf("envmerge: unknown strategy %d", p.Strategy)
	}
	return nil
}

// Merge combines sources in order according to the policy.
// Each source is a map of key→value pairs.
func (m *Merger) Merge(sources ...map[string]string) (map[string]string, error) {
	out := make(map[string]string)
	for _, src := range sources {
		for k, v := range src {
			existing, exists := out[k]
			switch {
			case !exists:
				out[k] = v
			case m.policy.Strategy == StrategyLastWins:
				out[k] = v
			case m.policy.Strategy == StrategyFirstWins:
				// keep existing
			case m.policy.Strategy == StrategyError:
				if existing != v {
					return nil, fmt.Errorf("envmerge: conflict on key %q", k)
				}
			}
		}
	}
	return out, nil
}
