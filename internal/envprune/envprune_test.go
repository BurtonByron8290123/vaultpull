package envprune

import (
	"testing"
)

func TestPruneRemovesStaleKey(t *testing.T) {
	pr := New(DefaultPolicy())
	current := map[string]string{"A": "1", "B": "2", "C": "3"}
	upstream := map[string]string{"A": "1", "C": "3"}
	out, res := pr.Apply(current, upstream)
	if _, ok := out["B"]; ok {
		t.Error("expected B to be removed")
	}
	if len(res.Removed) != 1 || res.Removed[0] != "B" {
		t.Errorf("unexpected removed list: %v", res.Removed)
	}
}

func TestPruneKeepsUpstreamKeys(t *testing.T) {
	pr := New(DefaultPolicy())
	current := map[string]string{"A": "1"}
	upstream := map[string]string{"A": "1"}
	out, res := pr.Apply(current, upstream)
	if _, ok := out["A"]; !ok {
		t.Error("expected A to be kept")
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removals, got %v", res.Removed)
	}
}

func TestPruneProtectedKeyNotRemoved(t *testing.T) {
	p := DefaultPolicy()
	p.ProtectedKeys = []string{"SECRET"}
	pr := New(p)
	current := map[string]string{"SECRET": "x", "OLD": "y"}
	upstream := map[string]string{}
	out, res := pr.Apply(current, upstream)
	if _, ok := out["SECRET"]; !ok {
		t.Error("expected SECRET to be protected")
	}
	if len(res.Protected) != 1 || res.Protected[0] != "SECRET" {
		t.Errorf("unexpected protected list: %v", res.Protected)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "OLD" {
		t.Errorf("unexpected removed list: %v", res.Removed)
	}
}

func TestPruneDryRunDoesNotMutate(t *testing.T) {
	p := DefaultPolicy()
	p.DryRun = true
	pr := New(p)
	current := map[string]string{"A": "1", "B": "2"}
	upstream := map[string]string{"A": "1"}
	out, res := pr.Apply(current, upstream)
	if _, ok := out["B"]; !ok {
		t.Error("dry-run should not remove B from returned map")
	}
	if !res.DryRun {
		t.Error("expected DryRun flag set in result")
	}
	if len(res.Removed) != 1 {
		t.Errorf("expected 1 reported removal, got %d", len(res.Removed))
	}
}

func TestSummaryFormat(t *testing.T) {
	r := Result{Removed: []string{"X"}, Protected: []string{"Y"}, DryRun: true}
	s := r.Summary()
	for _, want := range []string{"removed=1", "protected=1", "dry_run=true"} {
		if !containsStr(s, want) {
			t.Errorf("summary %q missing %q", s, want)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
