package tagmap_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/tagmap"
)

func TestApplyMatchingRule(t *testing.T) {
	m, err := tagmap.New([]tagmap.Rule{
		{Tag: "env", Value: "prod", Prefix: "PROD_"},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := m.Apply("DB_PASS", map[string]string{"env": "prod"})
	if got != "PROD_DB_PASS" {
		t.Fatalf("expected PROD_DB_PASS, got %s", got)
	}
}

func TestApplyNoMatchReturnsOriginal(t *testing.T) {
	m, _ := tagmap.New([]tagmap.Rule{
		{Tag: "env", Value: "prod", Prefix: "PROD_"},
	})
	got := m.Apply("DB_PASS", map[string]string{"env": "staging"})
	if got != "DB_PASS" {
		t.Fatalf("expected DB_PASS, got %s", got)
	}
}

func TestApplyEmptyPrefixReturnsKey(t *testing.T) {
	m, _ := tagmap.New([]tagmap.Rule{
		{Tag: "env", Value: "dev", Prefix: ""},
	})
	got := m.Apply("SECRET", map[string]string{"env": "dev"})
	if got != "SECRET" {
		t.Fatalf("expected SECRET, got %s", got)
	}
}

func TestApplyMapRewrites(t *testing.T) {
	m, _ := tagmap.New([]tagmap.Rule{
		{Tag: "tier", Value: "premium", Prefix: "PRE_"},
	})
	secrets := map[string]string{"API_KEY": "abc", "TOKEN": "xyz"}
	tags := map[string]string{"tier": "premium"}
	out := m.ApplyMap(secrets, tags)
	if out["PRE_API_KEY"] != "abc" {
		t.Errorf("expected PRE_API_KEY=abc, got %v", out)
	}
	if out["PRE_TOKEN"] != "xyz" {
		t.Errorf("expected PRE_TOKEN=xyz, got %v", out)
	}
}

func TestNewRejectsEmptyTag(t *testing.T) {
	_, err := tagmap.New([]tagmap.Rule{
		{Tag: "", Value: "prod", Prefix: "X_"},
	})
	if err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestNewRejectsEmptyValue(t *testing.T) {
	_, err := tagmap.New([]tagmap.Rule{
		{Tag: "env", Value: "", Prefix: "X_"},
	})
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestApplyFirstRuleWins(t *testing.T) {
	m, _ := tagmap.New([]tagmap.Rule{
		{Tag: "env", Value: "prod", Prefix: "FIRST_"},
		{Tag: "env", Value: "prod", Prefix: "SECOND_"},
	})
	got := m.Apply("KEY", map[string]string{"env": "prod"})
	if got != "FIRST_KEY" {
		t.Fatalf("expected FIRST_KEY, got %s", got)
	}
}
