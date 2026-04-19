package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/checkpoint"
)

func tempFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "checkpoint.json")
}

func TestGetMissingReturnsNotFound(t *testing.T) {
	s, err := checkpoint.NewStore(tempFile(t))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	_, ok := s.Get("secret/app")
	if ok {
		t.Fatal("expected not found for unknown path")
	}
}

func TestSetAndGet(t *testing.T) {
	s, err := checkpoint.NewStore(tempFile(t))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if err := s.Set("secret/app", "abc123"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	e, ok := s.Get("secret/app")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.SecretHash != "abc123" {
		t.Errorf("hash = %q, want abc123", e.SecretHash)
	}
	if e.PulledAt.IsZero() {
		t.Error("expected PulledAt to be set")
	}
	if time.Since(e.PulledAt) > 5*time.Second {
		t.Error("PulledAt is too far in the past")
	}
}

func TestPersistsAcrossReload(t *testing.T) {
	f := tempFile(t)
	s1, _ := checkpoint.NewStore(f)
	_ = s1.Set("secret/db", "hashXYZ")

	s2, err := checkpoint.NewStore(f)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e, ok := s2.Get("secret/db")
	if !ok {
		t.Fatal("expected persisted entry after reload")
	}
	if e.SecretHash != "hashXYZ" {
		t.Errorf("hash = %q, want hashXYZ", e.SecretHash)
	}
}

func TestFilePermissions(t *testing.T) {
	f := tempFile(t)
	s, _ := checkpoint.NewStore(f)
	_ = s.Set("secret/app", "h1")

	info, err := os.Stat(f)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("perm = %o, want 0600", perm)
	}
}
