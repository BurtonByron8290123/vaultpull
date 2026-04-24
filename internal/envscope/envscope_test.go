package envscope

import (
	"testing"
)

func newScoper(t *testing.T, p Policy) *Scoper {
	t.Helper()
	s, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestNoRulesPermitsAll(t *testing.T) {
	s := newScoper(t, Policy{})
	in := map[string]string{"APP_KEY": "a", "DB_PASS": "b"}
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestAllowPrefixFiltersKeys(t *testing.T) {
	s := newScoper(t, Policy{Allow: []string{"APP_"}})
	in := map[string]string{"APP_KEY": "a", "DB_PASS": "b", "APP_SECRET": "c"}
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should have been filtered out")
	}
}

func TestDenyPrefixRemovesKeys(t *testing.T) {
	s := newScoper(t, Policy{Deny: []string{"INTERNAL_"}})
	in := map[string]string{"APP_KEY": "a", "INTERNAL_SECRET": "b"}
	out := s.Apply(in)
	if _, ok := out["INTERNAL_SECRET"]; ok {
		t.Error("INTERNAL_SECRET should have been denied")
	}
	if _, ok := out["APP_KEY"]; !ok {
		t.Error("APP_KEY should be present")
	}
}

func TestDenyTakesPrecedenceOverAllow(t *testing.T) {
	s := newScoper(t, Policy{
		Allow: []string{"APP_"},
		Deny:  []string{"APP_INTERNAL_"},
	})
	in := map[string]string{"APP_KEY": "a", "APP_INTERNAL_TOKEN": "b"}
	out := s.Apply(in)
	if _, ok := out["APP_INTERNAL_TOKEN"]; ok {
		t.Error("APP_INTERNAL_TOKEN should have been denied")
	}
	if _, ok := out["APP_KEY"]; !ok {
		t.Error("APP_KEY should be allowed")
	}
}

func TestMultipleAllowPrefixes(t *testing.T) {
	s := newScoper(t, Policy{Allow: []string{"APP_", "SVC_"}})
	in := map[string]string{"APP_KEY": "a", "SVC_URL": "b", "DB_PASS": "c"}
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestNewRejectsBlankAllowPrefix(t *testing.T) {
	_, err := New(Policy{Allow: []string{"APP_", "  "}})
	if err == nil {
		t.Fatal("expected error for blank allow prefix")
	}
}

func TestNewRejectsBlankDenyPrefix(t *testing.T) {
	_, err := New(Policy{Deny: []string{"", "SVC_"}})
	if err == nil {
		t.Fatal("expected error for blank deny prefix")
	}
}
