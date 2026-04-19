package envclone

import (
	"os"
	"strings"
)

// FromEnv builds a Policy from environment variables.
//
//	VAULTPULL_CLONE_OVERWRITE=true
//	VAULTPULL_CLONE_DRY_RUN=true
func FromEnv() Policy {
	p := DefaultPolicy()
	if v := os.Getenv("VAULTPULL_CLONE_OVERWRITE"); strings.EqualFold(v, "true") || v == "1" {
		p.Overwrite = true
	}
	if v := os.Getenv("VAULTPULL_CLONE_DRY_RUN"); strings.EqualFold(v, "true") || v == "1" {
		p.DryRun = true
	}
	return p
}
