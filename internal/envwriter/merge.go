package envwriter

import (
	"fmt"
	"strings"
)

// MergeResult summarises the outcome of a merge operation.
type MergeResult struct {
	Added   []string
	Updated []string
	Unchanged []string
}

// String returns a human-readable summary of the merge result.
func (r MergeResult) String() string {
	var sb strings.Builder
	if len(r.Added) > 0 {
		fmt.Fprintf(&sb, "Added: %s\n", strings.Join(r.Added, ", "))
	}
	if len(r.Updated) > 0 {
		fmt.Fprintf(&sb, "Updated: %s\n", strings.Join(r.Updated, ", "))
	}
	if len(r.Unchanged) > 0 {
		fmt.Fprintf(&sb, "Unchanged: %s\n", strings.Join(r.Unchanged, ", "))
	}
	return sb.String()
}

// Merge combines existing and incoming secrets, returning the merged map
// and a MergeResult describing what changed.
func Merge(existing, incoming map[string]string) (map[string]string, MergeResult) {
	result := make(map[string]string, len(existing))
	for k, v := range existing {
		result[k] = v
	}

	var mr MergeResult
	for k, v := range incoming {
		old, exists := result[k]
		switch {
		case !exists:
			mr.Added = append(mr.Added, k)
		case old != v:
			mr.Updated = append(mr.Updated, k)
		default:
			mr.Unchanged = append(mr.Unchanged, k)
		}
		result[k] = v
	}

	return result, mr
}
