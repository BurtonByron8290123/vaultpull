package rollback_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultpull/internal/rollback"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "rollback-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestLatestReturnsNewestBackup(t *testing.T) {
	dir := tempDir(t)
	writeFile(t, filepath.Join(dir, ".env.20240101T120000.bak"), "old")
	writeFile(t, filepath.Join(dir, ".env.20240102T120000.bak"), "new")

	s, _ := rollback.New(dir)
	got, err := s.Latest(".env")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != ".env.20240102T120000.bak" {
		t.Errorf("expected newest backup, got %q", got)
	}
}

func TestLatestReturnsEmptyWhenNoBackups(t *testing.T) {
	dir := tempDir(t)
	s, _ := rollback.New(dir)
	got, err := s.Latest(".env")
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestRestoreWritesContentToDst(t *testing.T) {
	dir := tempDir(t)
	writeFile(t, filepath.Join(dir, ".env.20240101T120000.bak"), "SECRET=abc")

	s, _ := rollback.New(dir)
	dst := filepath.Join(dir, ".env")
	if err := s.Restore(".env", dst); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(dst)
	if string(data) != "SECRET=abc" {
		t.Errorf("unexpected content: %q", data)
	}
}

func TestRestoreErrorsWhenNoBackup(t *testing.T) {
	dir := tempDir(t)
	s, _ := rollback.New(dir)
	err := s.Restore(".env", filepath.Join(dir, ".env"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewRejectsEmptyDir(t *testing.T) {
	_, err := rollback.New("")
	if err == nil {
		t.Fatal("expected error for empty backupDir")
	}
}

func TestListReturnsSortedBackups(t *testing.T) {
	dir := tempDir(t)
	writeFile(t, filepath.Join(dir, ".env.20240103T000000.bak"), "c")
	writeFile(t, filepath.Join(dir, ".env.20240101T000000.bak"), "a")
	writeFile(t, filepath.Join(dir, ".env.20240102T000000.bak"), "b")

	s, _ := rollback.New(dir)
	list, err := s.List(".env")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(list))
	}
	if filepath.Base(list[0]) != ".env.20240101T000000.bak" {
		t.Errorf("expected oldest first, got %q", list[0])
	}
}
