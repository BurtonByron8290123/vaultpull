package envsort

import (
	"testing"
)

func newSorter(t *testing.T, p Policy) *Sorter {
	t.Helper()
	s, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestAlphabeticalOrder(t *testing.T) {
	s := newSorter(t, DefaultPolicy())
	m := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys := s.SortedKeys(m)
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestAlphabeticalDescOrder(t *testing.T) {
	s := newSorter(t, Policy{Strategy: AlphabeticalDesc, PrefixSep: "_"})
	m := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys := s.SortedKeys(m)
	want := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestPrefixGroupedOrder(t *testing.T) {
	s := newSorter(t, Policy{Strategy: PrefixGrouped, PrefixSep: "_"})
	m := map[string]string{
		"DB_HOST": "1",
		"APP_PORT": "2",
		"DB_PORT": "3",
		"APP_NAME": "4",
	}
	keys := s.SortedKeys(m)
	// APP group comes before DB group; within each group keys are sorted.
	want := []string{"APP_NAME", "APP_PORT", "DB_HOST", "DB_PORT"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	s := newSorter(t, DefaultPolicy())
	orig := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := s.Apply(orig)
	out["EXTRA"] = "val"
	if _, ok := orig["EXTRA"]; ok {
		t.Error("Apply mutated the input map")
	}
}

func TestDefaultPrefixSepFallback(t *testing.T) {
	// PrefixSep empty string should be normalised to "_".
	s, err := New(Policy{Strategy: PrefixGrouped, PrefixSep: ""})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if s.policy.PrefixSep != "_" {
		t.Errorf("expected default sep '_', got %q", s.policy.PrefixSep)
	}
}

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("VAULTPULL_SORT_STRATEGY", "")
	t.Setenv("VAULTPULL_SORT_PREFIX_SEP", "")
	p, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if p.Strategy != Alphabetical {
		t.Errorf("expected Alphabetical, got %v", p.Strategy)
	}
}

func TestFromEnvReadsStrategy(t *testing.T) {
	t.Setenv("VAULTPULL_SORT_STRATEGY", "prefix")
	p, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if p.Strategy != PrefixGrouped {
		t.Errorf("expected PrefixGrouped, got %v", p.Strategy)
	}
}

func TestFromEnvInvalidStrategyReturnsError(t *testing.T) {
	t.Setenv("VAULTPULL_SORT_STRATEGY", "random")
	_, err := FromEnv()
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}
