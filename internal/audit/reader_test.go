package audit_test

import (
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/audit"
)

func TestReadAllMissingFile(t *testing.T) {
	entries, err := audit.ReadAll("/nonexistent/audit.log")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestReadAllRoundTrip(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")
	l := audit.New(logPath)

	want := []audit.Entry{
		{Operation: "pull", Path: "secret/app", Target: ".env", Added: 2, Updated: 1},
		{Operation: "pull", Path: "secret/db", Target: ".env.db", Removed: 1, Error: "partial"},
	}
	for _, e := range want {
		if err := l.Record(e); err != nil {
			t.Fatalf("record: %v", err)
		}
	}

	got, err := audit.ReadAll(logPath)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d entries, got %d", len(want), len(got))
	}
	for i, e := range got {
		if e.Operation != want[i].Operation {
			t.Errorf("[%d] operation: want %q got %q", i, want[i].Operation, e.Operation)
		}
		if e.Path != want[i].Path {
			t.Errorf("[%d] path: want %q got %q", i, want[i].Path, e.Path)
		}
		if e.Added != want[i].Added {
			t.Errorf("[%d] added: want %d got %d", i, want[i].Added, e.Added)
		}
		if e.Error != want[i].Error {
			t.Errorf("[%d] error: want %q got %q", i, want[i].Error, e.Error)
		}
		if e.Timestamp.IsZero() {
			t.Errorf("[%d] expected non-zero timestamp", i)
		}
	}
}
