package envrestore_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envrestore"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "envrestore-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func writeBackup(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeBackup: %v", err)
	}
	return p
}

func newRestorer(t *testing.T, dir string, dryRun bool) *envrestore.Restorer {
	t.Helper()
	r, err := envrestore.New(envrestore.Policy{BackupDir: dir, DryRun: dryRun}, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return r
}

func TestLatestReturnsNewestBackup(t *testing.T) {
	dir := tempDir(t)
	writeBackup(t, dir, ".env.20240101T000000.bak", "A=1")
	writeBackup(t, dir, ".env.20240102T000000.bak", "A=2")
	r := newRestorer(t, dir, false)
	got, err := r.Latest(".env")
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if filepath.Base(got) != ".env.20240102T000000.bak" {
		t.Errorf("expected newest backup, got %s", got)
	}
}

func TestLatestReturnsErrNoBackupFound(t *testing.T) {
	dir := tempDir(t)
	r := newRestorer(t, dir, false)
	_, err := r.Latest(".env")
	if !errors.Is(err, envrestore.ErrNoBackupFound) {
		t.Fatalf("expected ErrNoBackupFound, got %v", err)
	}
}

func TestRestoreWritesFile(t *testing.T) {
	dir := tempDir(t)
	src := writeBackup(t, dir, ".env.20240101T000000.bak", "SECRET=hello")
	dst := filepath.Join(dir, ".env")
	r := newRestorer(t, dir, false)
	if err := r.Restore(dst, src); err != nil {
		t.Fatalf("Restore: %v", err)
	}
	data, _ := os.ReadFile(dst)
	if string(data) != "SECRET=hello" {
		t.Errorf("unexpected content: %s", data)
	}
}

func TestRestoreDryRunDoesNotWrite(t *testing.T) {
	dir := tempDir(t)
	src := writeBackup(t, dir, ".env.20240101T000000.bak", "SECRET=hello")
	dst := filepath.Join(dir, ".env")
	var buf bytes.Buffer
	r, _ := envrestore.New(envrestore.Policy{BackupDir: dir, DryRun: true}, &buf)
	if err := r.Restore(dst, src); err != nil {
		t.Fatalf("Restore dry-run: %v", err)
	}
	if _, err := os.Stat(dst); !errors.Is(err, os.ErrNotExist) {
		t.Error("expected dst not to be written in dry-run mode")
	}
	if buf.Len() == 0 {
		t.Error("expected dry-run output")
	}
}

func TestListBackupsSortedOldestFirst(t *testing.T) {
	dir := tempDir(t)
	writeBackup(t, dir, ".env.20240103T000000.bak", "")
	writeBackup(t, dir, ".env.20240101T000000.bak", "")
	writeBackup(t, dir, ".env.20240102T000000.bak", "")
	r := newRestorer(t, dir, false)
	list, err := r.ListBackups(".env")
	if err != nil {
		t.Fatalf("ListBackups: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 backups, got %d", len(list))
	}
	if filepath.Base(list[0]) != ".env.20240101T000000.bak" {
		t.Errorf("expected oldest first, got %s", list[0])
	}
}

func TestNewRejectsEmptyBackupDir(t *testing.T) {
	_, err := envrestore.New(envrestore.Policy{}, nil)
	if err == nil {
		t.Error("expected error for empty BackupDir")
	}
}
