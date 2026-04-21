package envrewrite

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyReplacesMatchingValue(t *testing.T) {
	rw, err := New(Policy{Rules: []Rule{
		{Pattern: `^http://`, Replacement: "https://"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	out := rw.Apply(map[string]string{"URL": "http://example.com"})
	if got := out["URL"]; got != "https://example.com" {
		t.Errorf("expected https://example.com, got %s", got)
	}
}

func TestApplyNoMatchLeavesValueUnchanged(t *testing.T) {
	rw, _ := New(Policy{Rules: []Rule{
		{Pattern: `^http://`, Replacement: "https://"},
	}})
	out := rw.Apply(map[string]string{"URL": "https://example.com"})
	if got := out["URL"]; got != "https://example.com" {
		t.Errorf("unexpected change: %s", got)
	}
}

func TestApplyKeyGlobScopesRule(t *testing.T) {
	rw, _ := New(Policy{Rules: []Rule{
		{Pattern: `localhost`, Replacement: "prod", KeyGlob: "DB_"},
	}})
	src := map[string]string{
		"DB_HOST":  "localhost",
		"API_HOST": "localhost",
	}
	out := rw.Apply(src)
	if out["DB_HOST"] != "prod" {
		t.Errorf("DB_HOST should be rewritten, got %s", out["DB_HOST"])
	}
	if out["API_HOST"] != "localhost" {
		t.Errorf("API_HOST should not be rewritten, got %s", out["API_HOST"])
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	rw, _ := New(Policy{Rules: []Rule{
		{Pattern: `foo`, Replacement: "bar"},
	}})
	src := map[string]string{"KEY": "foo"}
	rw.Apply(src)
	if src["KEY"] != "foo" {
		t.Error("input map was mutated")
	}
}

func TestApplyMultipleRulesAppliedInOrder(t *testing.T) {
	rw, _ := New(Policy{Rules: []Rule{
		{Pattern: `a`, Replacement: "b"},
		{Pattern: `b`, Replacement: "c"},
	}})
	out := rw.Apply(map[string]string{"K": "a"})
	if out["K"] != "c" {
		t.Errorf("expected c, got %s", out["K"])
	}
}

func TestNewRejectsInvalidPattern(t *testing.T) {
	_, err := New(Policy{Rules: []Rule{
		{Pattern: `[invalid`},
	}})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestNewRejectsEmptyPattern(t *testing.T) {
	_, err := New(Policy{Rules: []Rule{
		{Pattern: "", Replacement: "x"},
	}})
	if err == nil {
		t.Error("expected error for empty pattern")
	}
}

func TestLoadConfigReadsYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rewrite.yaml")
	content := "rules:\n  - pattern: \"foo\"\n    replacement: \"bar\"\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	p, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Rules) != 1 || p.Rules[0].Pattern != "foo" {
		t.Errorf("unexpected policy: %+v", p)
	}
}

func TestFromEnvNoopWhenEnvAbsent(t *testing.T) {
	t.Setenv("VAULTPULL_REWRITE_CONFIG", "")
	rw, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	out := rw.Apply(map[string]string{"K": "v"})
	if out["K"] != "v" {
		t.Error("expected noop rewriter")
	}
}
