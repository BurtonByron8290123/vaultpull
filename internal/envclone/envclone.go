// Package envclone copies secrets from one Vault path to another local env file.
package envclone

import (
	"fmt"
	"os"

	"github.com/example/vaultpull/internal/env"
)

// Policy controls clone behaviour.
type Policy struct {
	// Overwrite existing keys in the destination file.
	Overwrite bool
	// DryRun prints what would change without writing.
	DryRun bool
	// FileMode is the permission used when creating the destination file.
	FileMode os.FileMode
}

// DefaultPolicy returns sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		Overwrite: false,
		DryRun:    false,
		FileMode:  0600,
	}
}

// Cloner copies a map of secrets into a destination .env file.
type Cloner struct {
	policy Policy
}

// New returns a Cloner with the given policy.
func New(p Policy) (*Cloner, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	return &Cloner{policy: p}, nil
}

// Clone writes src secrets into dst, respecting the policy.
func (c *Cloner) Clone(dst string, src map[string]string) (int, error) {
	existing := map[string]string{}
	if entries, err := env.Parse(dst); err == nil {
		existing = env.ToMap(entries)
	}

	merged := make(map[string]string, len(existing))
	for k, v := range existing {
		merged[k] = v
	}

	written := 0
	for k, v := range src {
		if _, exists := merged[k]; exists && !c.policy.Overwrite {
			continue
		}
		merged[k] = v
		written++
	}

	if c.policy.DryRun {
		return written, nil
	}

	if err := env.WriteFile(dst, merged, c.policy.FileMode); err != nil {
		return 0, fmt.Errorf("envclone: write %s: %w", dst, err)
	}
	return written, nil
}

func validate(p Policy) error {
	if p.FileMode == 0 {
		return fmt.Errorf("envclone: file mode must not be zero")
	}
	return nil
}
