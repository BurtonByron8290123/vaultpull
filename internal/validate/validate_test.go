package validate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePassesWhenAllRulesMet(t *testing.T) {
	v := New(Policy{Rules: []Rule{{Key: "DB_PASS", MinLen: 4}}})
	if err := v.Validate(map[string]string{"DB_PASS": "secret"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateMissingRequiredKey(t *testing.T) {
	v := New(Policy{Rules: []Rule{{Key: "API_KEY"}}})
	if err := v.Validate(map[string]string{}); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestValidateMinLenViolation(t *testing.T) {
	v := New(Policy{Rules: []Rule{{Key: "TOKEN", MinLen: 10}}})
	if err := v.Validate(map[string]string{"TOKEN": "short"}); err == nil {
		t.Fatal("expected error for short value")
	}
}

func TestValidateDisallowedSubstring(t *testing.T) {
	v := New(Policy{Rules: []Rule{{Key: "PASS", Disallow: []string{"changeme"}}})
	if err := v.Validate(map[string]string{"PASS": "changeme123"}); err == nil {
		t.Fatal("expected error for disallowed substring")
	}
}

func TestValidateMultipleViolationsCombined(t *testing.T) {
	v := New(Policy{Rules: []Rule{
		{Key: "A"},
		{Key: "B"},
	}})
	err := v.Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected combined error")
	}
}

func TestFromEnvReadsRequire(t *testing.T) {
	t.Setenv("VAULTPULL_VALIDATE_REQUIRE", "FOO, BAR")
	t.Setenv("VAULTPULL_VALIDATE_MIN_LEN", "5")
	p := FromEnv()
	if len(p.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(p.Rules))
	}
	if p.Rules[0].MinLen != 5 {
		t.Fatalf("expected MinLen 5, got %d", p.Rules[0].MinLen)
	}
}

func TestFromFileLoadsRules(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.yaml")
	content := "rules:\n  - key: SECRET\n    min_len: 8\n    disallow:\n      - placeholder\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	p, err := FromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Rules) != 1 || p.Rules[0].Key != "SECRET" {
		t.Fatalf("unexpected rules: %+v", p.Rules)
	}
}

func TestFromFileMissingFileReturnsError(t *testing.T) {
	_, err := FromFile("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
