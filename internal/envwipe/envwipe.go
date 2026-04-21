// Package envwipe provides functionality to securely remove specific keys
// from an env map or file, optionally dry-running the operation.
package envwipe

import (
	"fmt"
	"os"
	"sort"

	"github.com/vaultpull/internal/env"
)

// Policy controls how the wipe operation behaves.
type Policy struct {
	// Keys is the explicit list of keys to remove.
	Keys []string
	// DryRun reports what would be removed without modifying the file.
	DryRun bool
}

// Result holds the outcome of a wipe operation.
type Result struct {
	Removed []string
	Skipped []string
}

// Summary returns a human-readable description of the result.
func (r Result) Summary() string {
	return fmt.Sprintf("removed %d key(s), skipped %d key(s)", len(r.Removed), len(r.Skipped))
}

// Wiper removes keys from env maps and files.
type Wiper struct {
	policy Policy
}

// New creates a Wiper with the given policy.
func New(p Policy) (*Wiper, error) {
	if len(p.Keys) == 0 {
		return nil, fmt.Errorf("envwipe: at least one key must be specified")
	}
	return &Wiper{policy: p}, nil
}

// Apply removes the configured keys from m and returns the result.
// The original map is not mutated; a new map is returned.
func (w *Wiper) Apply(m map[string]string) (map[string]string, Result) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}

	var res Result
	for _, key := range w.policy.Keys {
		if _, ok := out[key]; ok {
			res.Removed = append(res.Removed, key)
			if !w.policy.DryRun {
				delete(out, key)
			}
		} else {
			res.Skipped = append(res.Skipped, key)
		}
	}
	sort.Strings(res.Removed)
	sort.Strings(res.Skipped)
	return out, res
}

// WipeFile reads the env file at path, removes the configured keys, and
// writes the result back. In dry-run mode the file is not modified.
func (w *Wiper) WipeFile(path string) (Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Result{}, fmt.Errorf("envwipe: read %s: %w", path, err)
	}

	entries, err := env.Parse(string(data))
	if err != nil {
		return Result{}, fmt.Errorf("envwipe: parse %s: %w", path, err)
	}

	wiped, res := w.Apply(entries)
	if w.policy.DryRun {
		return res, nil
	}

	if err := env.WriteFile(path, wiped, 0o600); err != nil {
		return Result{}, fmt.Errorf("envwipe: write %s: %w", path, err)
	}
	return res, nil
}
