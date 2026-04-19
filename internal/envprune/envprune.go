// Package envprune removes stale keys from a .env file that are no longer
// present in the upstream Vault secret set.
package envprune

import (
	"fmt"
	"sort"
)

// Policy controls pruning behaviour.
type Policy struct {
	// DryRun reports what would be removed without modifying anything.
	DryRun bool
	// ProtectedKeys are never removed regardless of upstream state.
	ProtectedKeys []string
}

// DefaultPolicy returns a Policy with safe defaults.
func DefaultPolicy() Policy {
	return Policy{DryRun: false}
}

// Result holds the outcome of a prune operation.
type Result struct {
	Removed []string
	Protected []string
	DryRun bool
}

// Summary returns a human-readable summary line.
func (r Result) Summary() string {
	return fmt.Sprintf("removed=%d protected=%d dry_run=%v", len(r.Removed), len(r.Protected), r.DryRun)
}

// Pruner removes stale env keys.
type Pruner struct {
	policy Policy
	protected map[string]struct{}
}

// New creates a Pruner with the given policy.
func New(p Policy) *Pruner {
	pm := make(map[string]struct{}, len(p.ProtectedKeys))
	for _, k := range p.ProtectedKeys {
		pm[k] = struct{}{}
	}
	return &Pruner{policy: p, protected: pm}
}

// Apply removes keys from current that are absent in upstream.
// It returns the pruned map and a Result describing what changed.
func (pr *Pruner) Apply(current, upstream map[string]string) (map[string]string, Result) {
	out := make(map[string]string, len(current))
	for k, v := range current {
		out[k] = v
	}

	var removed, protected []string
	for k := range current {
		if _, inUpstream := upstream[k]; inUpstream {
			continue
		}
		if _, isProtected := pr.protected[k]; isProtected {
			protected = append(protected, k)
			continue
		}
		removed = append(removed, k)
		if !pr.policy.DryRun {
			delete(out, k)
		}
	}
	sort.Strings(removed)
	sort.Strings(protected)
	return out, Result{Removed: removed, Protected: protected, DryRun: pr.policy.DryRun}
}
