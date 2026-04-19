package envpin_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envpin"
)

func tempFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "pins.json")
}

func sampleEntries() []envpin.Entry {
	return []envpin.Entry{
		{Path: "secret/app/db", Version: 3, PinnedAt: time.Unix(1700000000, 0).UTC(), Checksum: "abc123"},
		{Path: "secret/app/api", Version: 1, PinnedAt: time.Unix(1700000001, 0).UTC(), Checksum: "def456"},
	}
}

func TestSaveAndLoad(t *testing.T) {
	s := envpin.NewStore(tempFile(t))
	entries := sampleEntries()
	if err := s.Save(entries); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(loaded))
	}
	if loaded[0].Path != entries[0].Path || loaded[0].Checksum != entries[0].Checksum {
		t.Errorf("entry mismatch: got %+v", loaded[0])
	}
}

func TestLoadMissingFileReturnsErrPinNotFound(t *testing.T) {
	s := envpin.NewStore("/nonexistent/pins.json")
	_, err := s.Load()
	if err != envpin.ErrPinNotFound {
		t.Fatalf("expected ErrPinNotFound, got %v", err)
	}
}

func TestSaveFilePermissions(t *testing.T) {
	path := tempFile(t)
	s := envpin.NewStore(path)
	if err := s.Save(sampleEntries()); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected 0600, got %o", perm)
	}
}

func TestDiffDetectsDrift(t *testing.T) {
	pinned := []envpin.Entry{
		{Path: "secret/app/db", Checksum: "old"},
		{Path: "secret/app/api", Checksum: "same"},
	}
	current := []envpin.Entry{
		{Path: "secret/app/db", Checksum: "new"},
		{Path: "secret/app/api", Checksum: "same"},
	}
	drifted := envpin.Diff(pinned, current)
	if len(drifted) != 1 || drifted[0] != "secret/app/db" {
		t.Errorf("unexpected drift result: %v", drifted)
	}
}

func TestDiffNoDriftWhenChecksumsSame(t *testing.T) {
	pinned := []envpin.Entry{{Path: "secret/x", Checksum: "aaa"}}
	current := []envpin.Entry{{Path: "secret/x", Checksum: "aaa"}}
	if d := envpin.Diff(pinned, current); len(d) != 0 {
		t.Errorf("expected no drift, got %v", d)
	}
}

func TestDiffIgnoresNewPaths(t *testing.T) {
	pinned := []envpin.Entry{{Path: "secret/old", Checksum: "aaa"}}
	current := []envpin.Entry{
		{Path: "secret/old", Checksum: "aaa"},
		{Path: "secret/new", Checksum: "bbb"},
	}
	if d := envpin.Diff(pinned, current); len(d) != 0 {
		t.Errorf("expected no drift for new paths, got %v", d)
	}
}
