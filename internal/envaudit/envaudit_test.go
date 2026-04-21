package envaudit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envaudit"
)

func tempFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit.jsonl")
}

func fixedClock(ts time.Time) func() time.Time { return func() time.Time { return ts } }

func TestRecordAppendsEntries(t *testing.T) {
	path := tempFile(t)
	l := envaudit.New(envaudit.Policy{AuditPath: path, MaskValues: true})
	changes := map[string]envaudit.ChangeKind{
		"DB_PASS": envaudit.ChangeAdded,
		"API_KEY": envaudit.ChangeUpdated,
	}
	if err := l.Record(".env", changes); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, err := envaudit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}
}

func TestRecordNoopWhenPathEmpty(t *testing.T) {
	l := envaudit.New(envaudit.Policy{AuditPath: ""})
	if err := l.Record(".env", map[string]envaudit.ChangeKind{"X": envaudit.ChangeAdded}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecordNoopWhenChangesEmpty(t *testing.T) {
	path := tempFile(t)
	l := envaudit.New(envaudit.Policy{AuditPath: path})
	if err := l.Record(".env", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("file should not have been created")
	}
}

func TestRecordSetsFilePermissions(t *testing.T) {
	path := tempFile(t)
	l := envaudit.New(envaudit.Policy{AuditPath: path})
	_ = l.Record(".env", map[string]envaudit.ChangeKind{"K": envaudit.ChangeRemoved})
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Fatalf("want 0600, got %04o", perm)
	}
}

func TestReadAllMissingFileReturnsErrNoLog(t *testing.T) {
	_, err := envaudit.ReadAll("/no/such/audit.jsonl")
	if err != envaudit.ErrNoLog {
		t.Fatalf("want ErrNoLog, got %v", err)
	}
}

func TestReadAllRoundTrip(t *testing.T) {
	path := tempFile(t)
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	l := envaudit.New(envaudit.Policy{AuditPath: path, MaskValues: false})
	l.(*struct{ policy envaudit.Policy; clock func() time.Time }) // compile guard – use exported helper instead
	_ = ts
	_ = l
	// Record via public API and verify round-trip via ReadAll.
	l2 := envaudit.New(envaudit.Policy{AuditPath: path})
	_ = l2.Record("app.env", map[string]envaudit.ChangeKind{"TOKEN": envaudit.ChangeAdded})
	entries, err := envaudit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "TOKEN" {
		t.Fatalf("want key TOKEN, got %s", entries[0].Key)
	}
	if entries[0].Kind != envaudit.ChangeAdded {
		t.Fatalf("want kind added, got %s", entries[0].Kind)
	}
}
