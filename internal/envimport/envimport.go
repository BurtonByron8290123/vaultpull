// Package envimport reads an existing .env file and imports its key-value
// pairs into a map, optionally skipping keys that already exist in a
// supplied base map.
package envimport

import (
	"fmt"
	"os"

	"github.com/vaultpull/internal/env"
)

// Policy controls how conflicts between the imported file and the base map
// are resolved.
type Policy struct {
	// SkipExisting prevents imported keys from overwriting keys already
	// present in the base map.
	SkipExisting bool

	// IgnoreMissing suppresses the error when the source file does not exist.
	IgnoreMissing bool
}

// DefaultPolicy returns a Policy with safe defaults.
func DefaultPolicy() Policy {
	return Policy{
		SkipExisting:  true,
		IgnoreMissing: false,
	}
}

// Importer merges an .env file into a target map.
type Importer struct {
	policy Policy
}

// New returns an Importer configured with the given Policy.
func New(p Policy) *Importer {
	return &Importer{policy: p}
}

// Import reads src and merges its entries into base, returning the combined
// map. base is never mutated; a new map is always returned.
func (im *Importer) Import(src string, base map[string]string) (map[string]string, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) && im.policy.IgnoreMissing {
			result := make(map[string]string, len(base))
			for k, v := range base {
				result[k] = v
			}
			return result, nil
		}
		return nil, fmt.Errorf("envimport: read %s: %w", src, err)
	}

	parsed, err := env.Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("envimport: parse %s: %w", src, err)
	}

	result := make(map[string]string, len(base)+len(parsed))
	for k, v := range base {
		result[k] = v
	}

	for _, entry := range parsed {
		if _, exists := result[entry.Key]; exists && im.policy.SkipExisting {
			continue
		}
		result[entry.Key] = entry.Value
	}

	return result, nil
}
