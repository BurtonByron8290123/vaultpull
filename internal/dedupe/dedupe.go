// Package dedupe provides deduplication of secret keys across multiple Vault paths.
package dedupe

import "fmt"

// Policy controls how duplicate keys are handled.
type Policy struct {
	// LastWins causes the last seen value to overwrite earlier ones.
	LastWins bool
	// FailOnConflict returns an error when the same key appears with different values.
	FailOnConflict bool
}

// DefaultPolicy returns a Policy where the last value wins.
func DefaultPolicy() Policy {
	return Policy{LastWins: true}
}

// Deduper merges maps of secrets while applying the configured policy.
type Deduper struct {
	policy Policy
}

// New returns a Deduper using the given policy.
func New(p Policy) (*Deduper, error) {
	if p.LastWins && p.FailOnConflict {
		return nil, fmt.Errorf("dedupe: LastWins and FailOnConflict are mutually exclusive")
	}
	return &Deduper{policy: p}, nil
}

// Merge combines multiple secret maps into one, applying the deduplication policy.
// Maps are processed in order; later maps are considered "last".
func (d *Deduper) Merge(maps ...map[string]string) (map[string]string, error) {
	out := make(map[string]string)
	origin := make(map[string]string) // key -> first value seen

	for _, m := range maps {
		for k, v := range m {
			prev, exists := out[k]
			if !exists {
				out[k] = v
				origin[k] = v
				continue
			}
			if prev == v {
				continue
			}
			if d.policy.FailOnConflict {
				return nil, fmt.Errorf("dedupe: conflict on key %q: %q vs %q", k, origin[k], v)
			}
			if d.policy.LastWins {
				out[k] = v
			}
		}
	}
	return out, nil
}
