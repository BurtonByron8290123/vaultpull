package filter_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/filter"
)

var baseSecrets = map[string]string{
	"APP_SECRET":   "abc",
	"APP_KEY":      "def",
	"DB_PASSWORD":  "pass",
	"DB_HOST":      "localhost",
	"INTERNAL_KEY": "secret",
}

func TestNoRulesReturnsAll(t *testing.T) {
	f := filter.New(nil)
	result := f.Apply(baseSecrets)
	if len(result) != len(baseSecrets) {
		t.Errorf("expected %d keys, got %d", len(baseSecrets), len(result))
	}
}

func TestIncludePrefixFiltersKeys(t *testing.T) {
	f := filter.New([]string{"APP_"})
	result := f.Apply(baseSecrets)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_SECRET"]; !ok {
		t.Error("expected APP_SECRET in result")
	}
	if _, ok := result["APP_KEY"]; !ok {
		t.Error("expected APP_KEY in result")
	}
}

func TestExcludePrefixRemovesKeys(t *testing.T) {
	f := filter.New([]string{"!INTERNAL_"})
	result := f.Apply(baseSecrets)
	if _, ok := result["INTERNAL_KEY"]; ok {
		t.Error("expected INTERNAL_KEY to be excluded")
	}
	if len(result) != len(baseSecrets)-1 {
		t.Errorf("expected %d keys, got %d", len(baseSecrets)-1, len(result))
	}
}

func TestExcludeTakesPrecedenceOverInclude(t *testing.T) {
	f := filter.New([]string{"APP_", "!APP_KEY"})
	result := f.Apply(baseSecrets)
	if _, ok := result["APP_KEY"]; ok {
		t.Error("expected APP_KEY to be excluded")
	}
	if _, ok := result["APP_SECRET"]; !ok {
		t.Error("expected APP_SECRET to be included")
	}
}

func TestMultipleIncludePrefixes(t *testing.T) {
	f := filter.New([]string{"APP_", "DB_"})
	result := f.Apply(baseSecrets)
	if len(result) != 4 {
		t.Errorf("expected 4 keys, got %d", len(result))
	}
}

func TestEmptyPatternsIgnored(t *testing.T) {
	f := filter.New([]string{"", "  ", "APP_"})
	result := f.Apply(baseSecrets)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}
