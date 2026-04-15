package diff

import "fmt"

// ChangeType represents the kind of change for a secret key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Updated  ChangeType = "updated"
	Removed  ChangeType = "removed"
	Unchanged ChangeType = "unchanged"
)

// Change describes a single key-level change between two env maps.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// Result holds the full diff between two env snapshots.
type Result struct {
	Changes []Change
}

// HasChanges returns true if any meaningful change exists.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary string.
func (r *Result) Summary() string {
	var added, updated, removed int
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			added++
		case Updated:
			updated++
		case Removed:
			removed++
		}
	}
	return fmt.Sprintf("+%d added, ~%d updated, -%d removed", added, updated, removed)
}

// Compare computes the diff between an existing env map and incoming secrets.
// Values are masked; only presence/change is tracked.
func Compare(existing, incoming map[string]string) *Result {
	result := &Result{}

	for k, newVal := range incoming {
		if oldVal, ok := existing[k]; ok {
			if oldVal == newVal {
				result.Changes = append(result.Changes, Change{Key: k, Type: Unchanged})
			} else {
				result.Changes = append(result.Changes, Change{Key: k, Type: Updated, OldVal: mask(oldVal), NewVal: mask(newVal)})
			}
		} else {
			result.Changes = append(result.Changes, Change{Key: k, Type: Added, NewVal: mask(newVal)})
		}
	}

	for k := range existing {
		if _, ok := incoming[k]; !ok {
			result.Changes = append(result.Changes, Change{Key: k, Type: Removed, OldVal: mask(existing[k])})
		}
	}

	return result
}

// mask replaces a secret value with asterisks, revealing only the length.
func mask(val string) string {
	if len(val) == 0 {
		return ""
	}
	return fmt.Sprintf("[%d chars]", len(val))
}
