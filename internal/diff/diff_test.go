package diff

import (
	"strings"
	"testing"
)

func TestCompareAdded(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"NEW_KEY": "secret"}

	r := Compare(existing, incoming)

	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	if r.Changes[0].Type != Added {
		t.Errorf("expected Added, got %s", r.Changes[0].Type)
	}
	if r.Changes[0].Key != "NEW_KEY" {
		t.Errorf("expected key NEW_KEY, got %s", r.Changes[0].Key)
	}
}

func TestCompareUpdated(t *testing.T) {
	existing := map[string]string{"DB_PASS": "old"}
	incoming := map[string]string{"DB_PASS": "new"}

	r := Compare(existing, incoming)

	if r.Changes[0].Type != Updated {
		t.Errorf("expected Updated, got %s", r.Changes[0].Type)
	}
}

func TestCompareRemoved(t *testing.T) {
	existing := map[string]string{"OLD_KEY": "val"}
	incoming := map[string]string{}

	r := Compare(existing, incoming)

	if r.Changes[0].Type != Removed {
		t.Errorf("expected Removed, got %s", r.Changes[0].Type)
	}
}

func TestCompareUnchanged(t *testing.T) {
	existing := map[string]string{"KEY": "same"}
	incoming := map[string]string{"KEY": "same"}

	r := Compare(existing, incoming)

	if r.Changes[0].Type != Unchanged {
		t.Errorf("expected Unchanged, got %s", r.Changes[0].Type)
	}
	if r.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestSummaryFormat(t *testing.T) {
	existing := map[string]string{"OLD": "x", "SAME": "y"}
	incoming := map[string]string{"SAME": "y", "NEW": "z"}

	r := Compare(existing, incoming)
	summary := r.Summary()

	if !strings.Contains(summary, "+1 added") {
		t.Errorf("expected +1 added in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "-1 removed") {
		t.Errorf("expected -1 removed in summary, got: %s", summary)
	}
}

func TestMaskHidesValue(t *testing.T) {
	existing := map[string]string{"SECRET": "hunter2"}
	incoming := map[string]string{"SECRET": "newpassword"}

	r := Compare(existing, incoming)

	if r.Changes[0].OldVal == "hunter2" {
		t.Error("old value should be masked")
	}
	if r.Changes[0].NewVal == "newpassword" {
		t.Error("new value should be masked")
	}
}
