package depcheck_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/depcheck"
)

func TestCheckAllPresent(t *testing.T) {
	c := depcheck.New(depcheck.Policy{Required: []string{"A", "B"}})
	res := c.Check(map[string]string{"A": "1", "B": "2", "C": "3"})
	if !res.OK() {
		t.Fatalf("expected OK, got missing: %v", res.Missing)
	}
	if len(res.Present) != 2 {
		t.Fatalf("expected 2 present, got %d", len(res.Present))
	}
}

func TestCheckMissingKeys(t *testing.T) {
	c := depcheck.New(depcheck.Policy{Required: []string{"A", "MISSING"}})
	res := c.Check(map[string]string{"A": "1"})
	if res.OK() {
		t.Fatal("expected not OK")
	}
	if len(res.Missing) != 1 || res.Missing[0] != "MISSING" {
		t.Fatalf("unexpected missing: %v", res.Missing)
	}
}

func TestCheckErrNilWhenOK(t *testing.T) {
	c := depcheck.New(depcheck.Policy{Required: []string{"X"}})
	res := c.Check(map[string]string{"X": "val"})
	if res.Err() != nil {
		t.Fatalf("expected nil error, got %v", res.Err())
	}
}

func TestCheckErrNonNilWhenMissing(t *testing.T) {
	c := depcheck.New(depcheck.Policy{Required: []string{"X", "Y"}})
	res := c.Check(map[string]string{})
	if res.Err() == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestCheckEmptyPolicy(t *testing.T) {
	c := depcheck.New(depcheck.Policy{})
	res := c.Check(map[string]string{"A": "1"})
	if !res.OK() {
		t.Fatal("empty policy should always pass")
	}
}

func TestFromEnvParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_REQUIRED_KEYS", "FOO, BAR , BAZ")
	p := depcheck.FromEnv()
	if len(p.Required) != 3 {
		t.Fatalf("expected 3 keys, got %d: %v", len(p.Required), p.Required)
	}
}

func TestFromEnvEmptyVar(t *testing.T) {
	t.Setenv("VAULTPULL_REQUIRED_KEYS", "")
	p := depcheck.FromEnv()
	if len(p.Required) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(p.Required))
	}
}
