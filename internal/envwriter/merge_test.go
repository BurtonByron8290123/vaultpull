package envwriter

import (
	"testing"
)

func TestMergeAddsNewKeys(t *testing.T) {
	existing := map[string]string{"A": "1"}
	incoming := map[string]string{"B": "2"}

	merged, result := Merge(existing, incoming)

	if merged["B"] != "2" {
		t.Errorf("expected B=2, got %s", merged["B"])
	}
	if len(result.Added) != 1 || result.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", result.Added)
	}
}

func TestMergeUpdatesChangedKeys(t *testing.T) {
	existing := map[string]string{"A": "old"}
	incoming := map[string]string{"A": "new"}

	merged, result := Merge(existing, incoming)

	if merged["A"] != "new" {
		t.Errorf("expected A=new, got %s", merged["A"])
	}
	if len(result.Updated) != 1 || result.Updated[0] != "A" {
		t.Errorf("expected Updated=[A], got %v", result.Updated)
	}
}

func TestMergeUnchangedKeys(t *testing.T) {
	existing := map[string]string{"A": "same"}
	incoming := map[string]string{"A": "same"}

	_, result := Merge(existing, incoming)

	if len(result.Unchanged) != 1 || result.Unchanged[0] != "A" {
		t.Errorf("expected Unchanged=[A], got %v", result.Unchanged)
	}
	if len(result.Updated) != 0 {
		t.Errorf("expected no updates, got %v", result.Updated)
	}
}

func TestMergePreservesExistingKeysNotInIncoming(t *testing.T) {
	existing := map[string]string{"KEEP": "me", "UPDATE": "old"}
	incoming := map[string]string{"UPDATE": "new"}

	merged, _ := Merge(existing, incoming)

	if merged["KEEP"] != "me" {
		t.Errorf("expected KEEP=me to be preserved, got %s", merged["KEEP"])
	}
	if merged["UPDATE"] != "new" {
		t.Errorf("expected UPDATE=new, got %s", merged["UPDATE"])
	}
}

func TestMergeResultString(t *testing.T) {
	mr := MergeResult{
		Added:     []string{"NEW_KEY"},
		Updated:   []string{"CHANGED_KEY"},
		Unchanged: []string{"SAME_KEY"},
	}
	s := mr.String()
	if s == "" {
		t.Error("expected non-empty MergeResult string")
	}
}
