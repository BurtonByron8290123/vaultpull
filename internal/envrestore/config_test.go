package envrestore_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envrestore"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("VAULTPULL_BACKUP_DIR", "")
	t.Setenv("VAULTPULL_BACKUP_SUFFIX", "")
	t.Setenv("VAULTPULL_RESTORE_DRY_RUN", "")
	p := envrestore.FromEnv()
	if p.BackupDir != ".vaultpull/backups" {
		t.Errorf("unexpected BackupDir: %s", p.BackupDir)
	}
	if p.Suffix != ".bak" {
		t.Errorf("unexpected Suffix: %s", p.Suffix)
	}
	if p.DryRun {
		t.Error("expected DryRun false by default")
	}
}

func TestFromEnvReadsBackupDir(t *testing.T) {
	t.Setenv("VAULTPULL_BACKUP_DIR", "/tmp/mybackups")
	p := envrestore.FromEnv()
	if p.BackupDir != "/tmp/mybackups" {
		t.Errorf("unexpected BackupDir: %s", p.BackupDir)
	}
}

func TestFromEnvReadsSuffix(t *testing.T) {
	t.Setenv("VAULTPULL_BACKUP_SUFFIX", ".backup")
	p := envrestore.FromEnv()
	if p.Suffix != ".backup" {
		t.Errorf("unexpected Suffix: %s", p.Suffix)
	}
}

func TestFromEnvReadsDryRun(t *testing.T) {
	t.Setenv("VAULTPULL_RESTORE_DRY_RUN", "true")
	p := envrestore.FromEnv()
	if !p.DryRun {
		t.Error("expected DryRun true")
	}
}

func TestFromEnvIgnoresInvalidDryRun(t *testing.T) {
	t.Setenv("VAULTPULL_RESTORE_DRY_RUN", "notabool")
	p := envrestore.FromEnv()
	if p.DryRun {
		t.Error("expected DryRun false for invalid value")
	}
}
