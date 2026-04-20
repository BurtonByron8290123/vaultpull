package envpatch

import (
	"fmt"
	"os"
	"strconv"
)

// Policy controls Patcher behaviour.
type Policy struct {
	// IgnoreExisting skips OpSet when the key already exists in the base map.
	IgnoreExisting bool
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		IgnoreExisting: false,
	}
}

func (p Policy) validate() error {
	// currently no numeric bounds to check; reserved for future constraints
	return nil
}

// FromEnv builds a Policy from environment variables.
//
//	VAULTPULL_PATCH_IGNORE_EXISTING=true
func FromEnv() (Policy, error) {
	p := DefaultPolicy()

	if raw := os.Getenv("VAULTPULL_PATCH_IGNORE_EXISTING"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return Policy{}, fmt.Errorf("envpatch: invalid VAULTPULL_PATCH_IGNORE_EXISTING: %w", err)
		}
		p.IgnoreExisting = v
	}

	return p, nil
}
