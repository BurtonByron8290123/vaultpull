package envbackup

import (
	"os"
	"strconv"
	"time"
)

const (
	envBackupDir     = "VAULTPULL_BACKUP_DIR"
	envMaxBackups    = "VAULTPULL_BACKUP_MAX"
	envBackupSuffix  = "VAULTPULL_BACKUP_SUFFIX"
	envBackupMaxAge  = "VAULTPULL_BACKUP_MAX_AGE_DAYS"
	envBackupDryRun  = "VAULTPULL_BACKUP_DRY_RUN"

	defaultBackupDir    = ".vaultpull/backups"
	defaultMaxBackups   = 10
	defaultBackupSuffix = ".bak"
	defaultMaxAgeDays   = 30
)

// FromEnv builds a Policy from environment variables, falling back to
// defaults for any variable that is absent or unparseable.
//
// Recognised variables:
//
//	VAULTPULL_BACKUP_DIR          – directory where backups are written
//	VAULTPULL_BACKUP_MAX          – maximum number of backups to retain
//	VAULTPULL_BACKUP_SUFFIX       – filename suffix appended to each backup
//	VAULTPULL_BACKUP_MAX_AGE_DAYS – backups older than this many days are pruned
//	VAULTPULL_BACKUP_DRY_RUN      – when "true" no files are written or deleted
func FromEnv() Policy {
	p := DefaultPolicy()

	if v := os.Getenv(envBackupDir); v != "" {
		p.Dir = v
	}

	if v := os.Getenv(envMaxBackups); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			p.MaxBackups = n
		}
	}

	if v := os.Getenv(envBackupSuffix); v != "" {
		p.Suffix = v
	}

	if v := os.Getenv(envBackupMaxAge); v != "" {
		if days, err := strconv.Atoi(v); err == nil && days > 0 {
			p.MaxAge = time.Duration(days) * 24 * time.Hour
		}
	}

	if v := os.Getenv(envBackupDryRun); v == "true" || v == "1" {
		p.DryRun = true
	}

	return p
}
