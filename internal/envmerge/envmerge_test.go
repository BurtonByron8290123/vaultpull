package envmerge

import "testing"

func newMerger(t *testing.T, s Strategy) *Merger {
	t.Helper()
	m, err := New(Policy{Strategy: s})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return m
}

func TestLastWinsOverwritesKey(t *testing.T) {
	m := newMerger(t, StrategyLastWins)
	out, err := m.Merge(
		map[string]string{"A": "1"},
		map[string]string{"A": "2"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "2" {
		t.Errorf("expected 2, got %s", out["A"])
	}
}

func TestFirstWinsKeepsOriginal(t *testing.T) {
	m := newMerger(t, StrategyFirstWins)
	out, err := m.Merge(
		map[string]string{"A": "1"},
		map[string]string{"A": "2"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Errorf("expected 1, got %s", out["A"])
	}
}

func TestErrorStrategyConflictReturnsError(t *testing.T) {
	m := newMerger(t, StrategyError)
	_, err := m.Merge(
		map[string]string{"A": "1"},
		map[string]string{"A": "2"},
	)
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestErrorStrategySameValueNoError(t *testing.T) {
	m := newMerger(t, StrategyError)
	out, err := m.Merge(
		map[string]string{"A": "1"},
		map[string]string{"A": "1"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Errorf("expected 1, got %s", out["A"])
	}
}

func TestMergeDisjointSources(t *testing.T) {
	m := newMerger(t, StrategyLastWins)
	out, err := m.Merge(
		map[string]string{"A": "1"},
		map[string]string{"B": "2"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Errorf("unexpected result: %v", out)
	}
}

func TestNewRejectsInvalidStrategy(t *testing.T) {
	_, err := New(Policy{Strategy: Strategy(99)})
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}
