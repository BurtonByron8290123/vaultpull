package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/snapshot"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(filepath.Join(dir, "snap.json"))

	secrets := map[string]string{"DB_PASS": "secret1", "API_KEY": "abc123"}
	snap := snapshot.Build("secret/myapp", secrets)

	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if loaded.Path != "secret/myapp" {
		t.Errorf("path: got %q, want %q", loaded.Path, "secret/myapp")
	}
	if len(loaded.Entries) != 2 {
		t.Errorf("entries: got %d, want 2", len(loaded.Entries))
	}
}

func TestLoadMissingFileReturnsNil(t *testing.T) {
	store := snapshot.NewStore("/nonexistent/path/snap.json")
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Error("expected nil snapshot for missing file")
	}
}

func TestSaveFilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	store := snapshot.NewStore(path)

	snap := snapshot.Build("secret/test", map[string]string{"KEY": "val"})
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("permissions: got %o, want 0600", perm)
	}
}

func TestBuildSortsEntriesByKey(t *testing.T) {
	secrets := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	snap := snapshot.Build("secret/app", secrets)

	if snap.Entries[0].Key != "A_KEY" {
		t.Errorf("first key: got %q, want A_KEY", snap.Entries[0].Key)
	}
	if snap.Entries[2].Key != "Z_KEY" {
		t.Errorf("last key: got %q, want Z_KEY", snap.Entries[2].Key)
	}
}
