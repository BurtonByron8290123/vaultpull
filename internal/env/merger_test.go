package env

import (
	"testing"
)

func TestMergeOverwriteReplacesExistingKeys(t *testing.T) {
	base := map[string]string{"A": "old", "B": "keep"}
	incoming := map[string]string{"A": "new", "C": "added"}
	out := Merge(base, incoming, StrategyOverwrite)
	if out["A"] != "new" {
		t.Errorf("expected A=new, got %s", out["A"])
	}
	if out["B"] != "keep" {
		t.Errorf("expected B=keep, got %s", out["B"])
	}
	if out["C"] != "added" {
		t.Errorf("expected C=added, got %s", out["C"])
	}
}

func TestMergePreserveKeepsExistingKeys(t *testing.T) {
	base := map[string]string{"A": "old", "B": "keep"}
	incoming := map[string]string{"A": "new", "C": "added"}
	out := Merge(base, incoming, StrategyPreserve)
	if out["A"] != "old" {
		t.Errorf("expected A=old, got %s", out["A"])
	}
	if out["C"] != "added" {
		t.Errorf("expected C=added, got %s", out["C"])
	}
}

func TestMergeEmptyBase(t *testing.T) {
	out := Merge(map[string]string{}, map[string]string{"X": "1"}, StrategyOverwrite)
	if out["X"] != "1" {
		t.Errorf("expected X=1, got %s", out["X"])
	}
}

func TestMergeEmptyIncoming(t *testing.T) {
	base := map[string]string{"A": "1"}
	out := Merge(base, map[string]string{}, StrategyOverwrite)
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %s", out["A"])
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestSubtractRemovesKeys(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	remove := map[string]string{"B": "anything"}
	out := Subtract(base, remove)
	if _, ok := out["B"]; ok {
		t.Error("expected B to be removed")
	}
	if out["A"] != "1" || out["C"] != "3" {
		t.Error("expected A and C to be preserved")
	}
}

func TestSubtractEmptyRemove(t *testing.T) {
	base := map[string]string{"A": "1"}
	out := Subtract(base, map[string]string{})
	if len(out) != len(base) {
		t.Errorf("expected %d keys, got %d", len(base), len(out))
	}
}
