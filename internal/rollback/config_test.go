package rollback_test

import (
	"testing"

	"github.com/vaultpull/internal/rollback"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("VAULTPULL_BACKUP_DIR", "")
	t.Setenv("VAULTPULL_MAX_BACKUPS", "")

	cfg := rollback.FromEnv()
	if cfg.BackupDir != ".vaultpull/backups" {
		t.Errorf("unexpected BackupDir: %q", cfg.BackupDir)
	}
	if cfg.MaxBackups != 5 {
		t.Errorf("unexpected MaxBackups: %d", cfg.MaxBackups)
	}
}

func TestFromEnvReadsBackupDir(t *testing.T) {
	t.Setenv("VAULTPULL_BACKUP_DIR", "/tmp/mybackups")
	t.Setenv("VAULTPULL_MAX_BACKUPS", "")

	cfg := rollback.FromEnv()
	if cfg.BackupDir != "/tmp/mybackups" {
		t.Errorf("unexpected BackupDir: %q", cfg.BackupDir)
	}
}

func TestFromEnvReadsMaxBackups(t *testing.T) {
	t.Setenv("VAULTPULL_BACKUP_DIR", "")
	t.Setenv("VAULTPULL_MAX_BACKUPS", "10")

	cfg := rollback.FromEnv()
	if cfg.MaxBackups != 10 {
		t.Errorf("unexpected MaxBackups: %d", cfg.MaxBackups)
	}
}

func TestFromEnvIgnoresInvalidMaxBackups(t *testing.T) {
	t.Setenv("VAULTPULL_MAX_BACKUPS", "notanumber")

	cfg := rollback.FromEnv()
	if cfg.MaxBackups != 5 {
		t.Errorf("expected default 5, got %d", cfg.MaxBackups)
	}
}

func TestFromEnvZeroMaxBackupsIsAllowed(t *testing.T) {
	t.Setenv("VAULTPULL_MAX_BACKUPS", "0")

	cfg := rollback.FromEnv()
	if cfg.MaxBackups != 0 {
		t.Errorf("expected 0, got %d", cfg.MaxBackups)
	}
}
