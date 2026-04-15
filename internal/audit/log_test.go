package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/audit"
)

func TestRecordWritesEntry(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	l := audit.New(logPath)
	err := l.Record(audit.Entry{
		Operation: "pull",
		Path:      "secret/myapp",
		Target:    ".env",
		Added:     3,
		Updated:   1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, _ := os.Open(logPath)
	defer f.Close()

	var entry audit.Entry
	if err := json.NewDecoder(f).Decode(&entry); err != nil {
		t.Fatalf("decode entry: %v", err)
	}
	if entry.Operation != "pull" {
		t.Errorf("expected operation=pull, got %q", entry.Operation)
	}
	if entry.Added != 3 {
		t.Errorf("expected added=3, got %d", entry.Added)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecordAppendsMultipleEntries(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")
	l := audit.New(logPath)

	for i := 0; i < 3; i++ {
		if err := l.Record(audit.Entry{Operation: "pull", Path: "secret/app"}); err != nil {
			t.Fatalf("record %d: %v", i, err)
		}
	}

	f, _ := os.Open(logPath)
	defer f.Close()

	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() != "" {
			count++
		}
	}
	if count != 3 {
		t.Errorf("expected 3 log lines, got %d", count)
	}
}

func TestRecordNoopWhenPathEmpty(t *testing.T) {
	l := audit.New("")
	if err := l.Record(audit.Entry{Operation: "pull"}); err != nil {
		t.Errorf("expected no error for empty path, got %v", err)
	}
}

func TestRecordSetsFilePermissions(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")
	l := audit.New(logPath)
	_ = l.Record(audit.Entry{Operation: "pull"})

	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected perm 0600, got %o", perm)
	}
}
