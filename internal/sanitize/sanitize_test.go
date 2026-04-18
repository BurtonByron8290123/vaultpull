package sanitize_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/sanitize"
)

func TestNormalizeKeysUppercases(t *testing.T) {
	s := sanitize.New(sanitize.Policy{NormalizeKeys: true})
	out := s.Apply(map[string]string{"db_host": "localhost"})
	if _, ok := out["DB_HOST"]; !ok {
		t.Fatal("expected DB_HOST key")
	}
}

func TestStripInvalidKeysRemovesBadKeys(t *testing.T) {
	s := sanitize.New(sanitize.Policy{StripInvalidKeys: true})
	out := s.Apply(map[string]string{"1INVALID": "val", "VALID_KEY": "ok"})
	if _, ok := out["1INVALID"]; ok {
		t.Fatal("expected 1INVALID to be stripped")
	}
	if _, ok := out["VALID_KEY"]; !ok {
		t.Fatal("expected VALID_KEY to be kept")
	}
}

func TestTrimValuesRemovesWhitespace(t *testing.T) {
	s := sanitize.New(sanitize.Policy{TrimValues: true})
	out := s.Apply(map[string]string{"KEY": "  hello  "})
	if out["KEY"] != "hello" {
		t.Fatalf("expected 'hello', got %q", out["KEY"])
	}
}

func TestStripNullBytesRemovesNulls(t *testing.T) {
	s := sanitize.New(sanitize.Policy{StripNullBytes: true})
	out := s.Apply(map[string]string{"KEY": "val\x00ue"})
	if out["KEY"] != "value" {
		t.Fatalf("expected 'value', got %q", out["KEY"])
	}
}

func TestDefaultPolicyAppliesAll(t *testing.T) {
	s := sanitize.New(sanitize.DefaultPolicy())
	in := map[string]string{
		"lower_key": "  trimmed  ",
		"2bad":       "dropped",
		"good_key":   "val\x00ue",
	}
	out := s.Apply(in)
	if _, ok := out["2bad"]; ok {
		t.Fatal("2bad should be stripped")
	}
	if out["LOWER_KEY"] != "trimmed" {
		t.Fatalf("unexpected value: %q", out["LOWER_KEY"])
	}
	if out["GOOD_KEY"] != "value" {
		t.Fatalf("unexpected value: %q", out["GOOD_KEY"])
	}
}

func TestIsValidKey(t *testing.T) {
	cases := []struct {
		key   string
		valid bool
	}{
		{"VALID", true},
		{"_ALSO_VALID", true},
		{"valid123", true},
		{"1INVALID", false},
		{"has-dash", false},
		{"", false},
	}
	for _, c := range cases {
		got := sanitize.IsValidKey(c.key)
		if got != c.valid {
			t.Errorf("IsValidKey(%q) = %v, want %v", c.key, got, c.valid)
		}
	}
}
