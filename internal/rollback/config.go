package rollback

import (
	"os"
	"strconv"
)

// Config holds rollback configuration.
type Config struct {
	// BackupDir is the directory where rotation backups are stored.
	BackupDir string
	// MaxBackups is the maximum number of backups to consider (0 = unlimited).
	MaxBackups int
}

// FromEnv builds a Config from environment variables with sensible defaults.
//
//	VAULTPULL_BACKUP_DIR   – directory for backups (default: ".vaultpull/backups")
//	VAULTPULL_MAX_BACKUPS  – max backups to keep  (default: 5)
func FromEnv() Config {
	dir := os.Getenv("VAULTPULL_BACKUP_DIR")
	if dir == "" {
		dir = ".vaultpull/backups"
	}

	max := 5
	if v := os.Getenv("VAULTPULL_MAX_BACKUPS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			max = n
		}
	}

	return Config{
		BackupDir:  dir,
		MaxBackups: max,
	}
}
