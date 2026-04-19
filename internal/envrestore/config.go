package envrestore

import (
	"os"
	"strconv"
)

// FromEnv builds a Policy from environment variables.
//
//	VAULTPULL_BACKUP_DIR   – directory containing rotation backups
//	VAULTPULL_BACKUP_SUFFIX – file suffix used by the rotation package
//	VAULTPULL_RESTORE_DRY_RUN – set to "true" to enable dry-run mode
func FromEnv() Policy {
	p := Policy{
		BackupDir: ".vaultpull/backups",
		Suffix:    ".bak",
	}
	if v := os.Getenv("VAULTPULL_BACKUP_DIR"); v != "" {
		p.BackupDir = v
	}
	if v := os.Getenv("VAULTPULL_BACKUP_SUFFIX"); v != "" {
		p.Suffix = v
	}
	if v := os.Getenv("VAULTPULL_RESTORE_DRY_RUN"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			p.DryRun = b
		}
	}
	return p
}
