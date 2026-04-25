// Package envdiff compares two env file snapshots and produces a human-readable
// change report suitable for dry-run and audit output.
package envdiff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ChangeKind describes the type of change for a single key.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Updated ChangeKind = "updated"
	Removed ChangeKind = "removed"
)

// Change represents a single key-level diff entry.
type Change struct {
	Key  string
	Kind ChangeKind
	Old  string
	New  string
}

// Report holds the full diff between two env maps.
type Report struct {
	Changes []Change
}

// HasChanges returns true when at least one change is present.
func (r *Report) HasChanges() bool { return len(r.Changes) > 0 }

// Summary returns a one-line human-readable summary.
func (r *Report) Summary() string {
	var added, updated, removed int
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			added++
		case Updated:
			updated++
		case Removed:
			removed++
		}
	}
	return fmt.Sprintf("+%d added  ~%d updated  -%d removed", added, updated, removed)
}

// FilterByKind returns a new Report containing only changes of the given kind.
func (r *Report) FilterByKind(kind ChangeKind) *Report {
	var filtered []Change
	for _, c := range r.Changes {
		if c.Kind == kind {
			filtered = append(filtered, c)
		}
	}
	return &Report{Changes: filtered}
}

// Compare produces a Report by diffing prev and next env maps.
func Compare(prev, next map[string]string) *Report {
	seen := make(map[string]bool)
	var changes []Change

	for k, nv := range next {
		seen[k] = true
		if ov, ok := prev[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Added, New: nv})
		} else if ov != nv {
			changes = append(changes, Change{Key: k, Kind: Updated, Old: ov, New: nv})
		}
	}
	for k, ov := range prev {
		if !seen[k] {
			changes = append(changes, Change{Key: k, Kind: Removed, Old: ov})
		}
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Key < changes[j].Key })
	return &Report{Changes: changes}
}

// Fprint writes a colourless textual diff to w.
func Fprint(w io.Writer, r *Report) {
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			fmt.Fprintf(w, "+ %s=%s\n", c.Key, mask(c.New))
		case Removed:
			fmt.Fprintf(w, "- %s\n", c.Key)
		case Updated:
			fmt.Fprintf(w, "~ %s=%s -> %s\n", c.Key, mask(c.Old), mask(c.New))
		}
	}
}

func mask(v string) string {
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return strings.Repeat("*", len(v)-2) + v[len(v)-2:]
}
