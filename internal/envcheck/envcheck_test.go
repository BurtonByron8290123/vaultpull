package envcheck

import (
	"testing"
)

func newChecker(t *testing.T, p Policy) *Checker {
	t.Helper()
	c, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func TestCheckPassesWhenAllPresent(t *testing.T) {
	c := newChecker(t, Policy{Required: []string{"FOO", "BAR"}})
	env := map[string]string{"FOO": "1", "BAR": "2"}
	if err := c.Check(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckFailsWhenRequiredKeyMissing(t *testing.T) {
	c := newChecker(t, Policy{Required: []string{"FOO", "MISSING"}})
	env := map[string]string{"FOO": "1"}
	if err := c.Check(env); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCheckFailsWhenNonEmptyKeyIsEmpty(t *testing.T) {
	c := newChecker(t, Policy{NonEmpty: []string{"TOKEN"}})
	env := map[string]string{"TOKEN": "   "}
	if err := c.Check(env); err == nil {
		t.Fatal("expected error for blank value")
	}
}

func TestCheckPassesWhenNonEmptyKeyHasValue(t *testing.T) {
	c := newChecker(t, Policy{NonEmpty: []string{"TOKEN"}})
	env := map[string]string{"TOKEN": "secret"}
	if err := c.Check(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckNonEmptyKeyMissingReturnsError(t *testing.T) {
	c := newChecker(t, Policy{NonEmpty: []string{"MISSING"}})
	env := map[string]string{}
	if err := c.Check(env); err == nil {
		t.Fatal("expected error for missing non-empty key")
	}
}

func TestCheckEmptyPolicyAlwaysPasses(t *testing.T) {
	c := newChecker(t, Policy{})
	if err := c.Check(map[string]string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewRejectsBlankRequiredKey(t *testing.T) {
	_, err := New(Policy{Required: []string{"", "FOO"}})
	if err == nil {
		t.Fatal("expected error for blank required key")
	}
}

func TestNewRejectsBlankNonEmptyKey(t *testing.T) {
	_, err := New(Policy{NonEmpty: []string{"  "}})
	if err == nil {
		t.Fatal("expected error for blank non-empty key")
	}
}

func TestSplitCSVHandlesEmptyString(t *testing.T) {
	if got := splitCSV(""); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestSplitCSVTrimsSpaces(t *testing.T) {
	got := splitCSV(" A , B , C ")
	if len(got) != 3 || got[0] != "A" || got[1] != "B" || got[2] != "C" {
		t.Fatalf("unexpected result: %v", got)
	}
}
