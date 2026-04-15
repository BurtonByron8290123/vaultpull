package rotation

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRotateCreatesBackup(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")
	backupDir := filepath.Join(tmpDir, "backups")

	if err := os.WriteFile(envFile, []byte("SECRET=abc\n"), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	r := New(backupDir, 5)
	if err := r.Rotate(envFile); err != nil {
		t.Fatalf("Rotate() error: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(backupDir, ".env.*.bak"))
	if len(matches) != 1 {
		t.Errorf("expected 1 backup, got %d", len(matches))
	}
}

func TestRotateNoopWhenFileMissing(t *testing.T) {
	tmpDir := t.TempDir()
	r := New(filepath.Join(tmpDir, "backups"), 5)

	if err := r.Rotate(filepath.Join(tmpDir, ".env")); err != nil {
		t.Errorf("expected no error for missing file, got: %v", err)
	}
}

func TestPruneKeepsMaxBackups(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")
	backupDir := filepath.Join(tmpDir, "backups")

	if err := os.MkdirAll(backupDir, 0700); err != nil {
		t.Fatalf("setup: %v", err)
	}

	for i := 1; i <= 6; i++ {
		name := filepath.Join(backupDir, fmt.Sprintf(".env.2024010%dT120000.bak", i))
		_ = os.WriteFile(name, []byte("x"), 0600)
	}

	if err := os.WriteFile(envFile, []byte("SECRET=new\n"), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	r := New(backupDir, 5)
	if err := r.Rotate(envFile); err != nil {
		t.Fatalf("Rotate() error: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(backupDir, ".env.*.bak"))
	if len(matches) != 5 {
		t.Errorf("expected 5 backups after pruning, got %d", len(matches))
	}
}
