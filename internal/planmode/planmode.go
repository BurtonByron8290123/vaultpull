// Package planmode provides a dry-run planner that shows what changes would
// be applied without writing any files.
package planmode

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultpull/internal/diff"
)

// Change describes a single planned key operation.
type Change struct {
	Key    string
	OldVal string
	NewVal string
	Op     diff.ChangeType
}

// Plan holds the full set of planned changes for a target env file.
type Plan struct {
	Path    string
	Changes []Change
}

// HasChanges returns true when at least one non-unchanged entry exists.
func (p *Plan) HasChanges() bool {
	for _, c := range p.Changes {
		if c.Op != diff.Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a short human-readable summary line.
func (p *Plan) Summary() string {
	var added, updated, removed int
	for _, c := range p.Changes {
		switch c.Op {
		case diff.Added:
			added++
		case diff.Updated:
			updated++
		case diff.Removed:
			removed++
		}
	}
	return fmt.Sprintf("%s: +%d ~%d -%d", p.Path, added, updated, removed)
}

// Build constructs a Plan by comparing existing env values against incoming ones.
func Build(path string, existing, incoming map[string]string) *Plan {
	seen := make(map[string]bool)
	var changes []Change

	for k, newVal := range incoming {
		seen[k] = true
		if oldVal, ok := existing[k]; ok {
			if oldVal == newVal {
				changes = append(changes, Change{Key: k, OldVal: oldVal, NewVal: newVal, Op: diff.Unchanged})
			} else {
				changes = append(changes, Change{Key: k, OldVal: oldVal, NewVal: newVal, Op: diff.Updated})
			}
		} else {
			changes = append(changes, Change{Key: k, NewVal: newVal, Op: diff.Added})
		}
	}

	for k, oldVal := range existing {
		if !seen[k] {
			changes = append(changes, Change{Key: k, OldVal: oldVal, Op: diff.Removed})
		}
	}

	sort.Slice(changes, func(i, j int) bool { return changes[i].Key < changes[j].Key })
	return &Plan{Path: path, Changes: changes}
}

// Print writes the plan to w, defaulting to os.Stdout when w is nil.
func Print(p *Plan, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, "Plan for %s\n", p.Path)
	for _, c := range p.Changes {
		switch c.Op {
		case diff.Added:
			fmt.Fprintf(w, "  + %s\n", c.Key)
		case diff.Updated:
			fmt.Fprintf(w, "  ~ %s\n", c.Key)
		case diff.Removed:
			fmt.Fprintf(w, "  - %s\n", c.Key)
		}
	}
	if !p.HasChanges() {
		fmt.Fprintln(w, "  (no changes)")
	}
}
