package snapshot_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/snapshot"
)

func findChange(changes []snapshot.Change, key string) *snapshot.Change {
	for i := range changes {
		if changes[i].Key == key {
			return &changes[i]
		}
	}
	return nil
}

func TestDiffNilPrevAllAdded(t *testing.T) {
	curr := snapshot.Build("secret/app", map[string]string{"A": "1", "B": "2"})
	changes := snapshot.Diff(nil, curr)

	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	for _, c := range changes {
		if c.Kind != snapshot.ChangeAdded {
			t.Errorf("key %q: expected added, got %q", c.Key, c.Kind)
		}
	}
}

func TestDiffDetectsUpdated(t *testing.T) {
	prev := snapshot.Build("secret/app", map[string]string{"KEY": "old"})
	curr := snapshot.Build("secret/app", map[string]string{"KEY": "new"})

	changes := snapshot.Diff(&prev, curr)
	c := findChange(changes, "KEY")
	if c == nil {
		t.Fatal("expected change for KEY")
	}
	if c.Kind != snapshot.ChangeUpdated {
		t.Errorf("kind: got %q, want updated", c.Kind)
	}
}

func TestDiffDetectsRemoved(t *testing.T) {
	prev := snapshot.Build("secret/app", map[string]string{"OLD": "val", "KEEP": "x"})
	curr := snapshot.Build("secret/app", map[string]string{"KEEP": "x"})

	changes := snapshot.Diff(&prev, curr)
	c := findChange(changes, "OLD")
	if c == nil {
		t.Fatal("expected change for OLD")
	}
	if c.Kind != snapshot.ChangeRemoved {
		t.Errorf("kind: got %q, want removed", c.Kind)
	}
}

func TestDiffUnchangedProducesNoChanges(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	prev := snapshot.Build("secret/app", secrets)
	curr := snapshot.Build("secret/app", secrets)

	changes := snapshot.Diff(&prev, curr)
	if len(changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(changes))
	}
}
