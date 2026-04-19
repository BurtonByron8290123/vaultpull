package envprune

import (
	"os"
	"strings"
)

// FromEnv builds a Policy from environment variables.
//
//	VAULTPULL_PRUNE_DRY_RUN        – "true" enables dry-run mode
//	VAULTPULL_PRUNE_PROTECTED_KEYS – comma-separated list of protected keys
func FromEnv() Policy {
	p := DefaultPolicy()

	if v := os.Getenv("VAULTPULL_PRUNE_DRY_RUN"); strings.EqualFold(v, "true") {
		p.DryRun = true
	}

	if v := os.Getenv("VAULTPULL_PRUNE_PROTECTED_KEYS"); v != "" {
		for _, k := range strings.Split(v, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				p.ProtectedKeys = append(p.ProtectedKeys, k)
			}
		}
	}

	return p
}
