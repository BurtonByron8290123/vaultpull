// Package envcompare compares two .env files and reports differences.
package envcompare

import (
	"fmt"
	"io"
	"sort"

	"github.com/example/vaultpull/internal/env"
)

// ChangeKind describes the type of difference between two env files.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Updated ChangeKind = "updated"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result the full comparison output.
type Result struct {
	Changes []Change
}

// Summary returns a human-readable one-line summary.
func (r *Result) Summary() string {
	var added, removed, updated int
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Updated:
			updated++
		}
	}
	return fmt.Sprintf("+%d added, ~%d updated, -%d removed", added, updated, removed)
}

// Compare returns the diff between src and dst env maps.
func Compare(src, dst map[string]string) *Result {
	changes := []Change{}
	for k, dv := range dst {
		if sv, ok := src[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Added, NewValue: dv})
		} else if sv != dv {
			changes = append(changes, Change{Key: k, Kind: Updated, OldValue: sv, NewValue: dv})
		}
	}
	for k, sv := range src {
		if _, ok := dst[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: sv})
		}
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Key < changes[j].Key })
	return &Result{Changes: changes}
}

// CompareFiles loads two .env files and compares them.
func CompareFiles(srcPath, dstPath string) (*Result, error) {
	src, err := loadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("envcompare: load src: %w", err)
	}
	dst, err := loadFile(dstPath)
	if err != nil {
		return nil, fmt.Errorf("envcompare: load dst: %w", err)
	}
	return Compare(src, dst), nil
}

func loadFile(path string) (map[string]string, error) {
	entries, err := env.Parse(path)
	if err != nil {
		return nil, err
	}
	return env.ToMap(entries), nil
}

// Fprint writes a human-readable diff to w.
func Fprint(w io.Writer, r *Result) {
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			fmt.Fprintf(w, "+ %s=%s\n", c.Key, c.NewValue)
		case Removed:
			fmt.Fprintf(w, "- %s=%s\n", c.Key, c.OldValue)
		case Updated:
			fmt.Fprintf(w, "~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
		}
	}
}
