package truncate_test

import (
	"strings"
	"testing"

	"github.com/example/vaultpull/internal/truncate"
)

func TestNoLimitReturnsOriginal(t *testing.T) {
	p := truncate.Policy{MaxLen: 0, Suffix: "..."}
	got := p.Apply(strings.Repeat("x", 200))
	if len(got) != 200 {
		t.Fatalf("expected 200, got %d", len(got))
	}
}

func TestShortValueUnchanged(t *testing.T) {
	p := truncate.Policy{MaxLen: 50, Suffix: "..."}
	const val = "hello"
	if got := p.Apply(val); got != val {
		t.Fatalf("expected %q, got %q", val, got)
	}
}

func TestTruncatesLongValue(t *testing.T) {
	p := truncate.Policy{MaxLen: 10, Suffix: "..."}
	got := p.Apply("abcdefghijklmnop")
	if len([]rune(got)) != 10 {
		t.Fatalf("expected length 10, got %d", len([]rune(got)))
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
}

func TestApplyMapTruncatesValues(t *testing.T) {
	p := truncate.Policy{MaxLen: 5, Suffix: "~"}
	m := map[string]string{
		"SHORT": "hi",
		"LONG":  "abcdefgh",
	}
	out := p.ApplyMap(m)
	if out["SHORT"] != "hi" {
		t.Errorf("SHORT should be unchanged")
	}
	if !strings.HasSuffix(out["LONG"], "~") {
		t.Errorf("LONG should be truncated, got %q", out["LONG"])
	}
	if len([]rune(out["LONG"])) != 5 {
		t.Errorf("LONG length should be 5, got %d", len([]rune(out["LONG"])))
	}
}

func TestValidateRejectsNegativeMaxLen(t *testing.T) {
	p := truncate.Policy{MaxLen: -1, Suffix: "..."}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for negative MaxLen")
	}
}

func TestValidateRejectsSuffixLongerThanMax(t *testing.T) {
	p := truncate.Policy{MaxLen: 3, Suffix: "..."}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error when suffix >= MaxLen")
	}
}

func TestDefaultPolicyHasNoLimit(t *testing.T) {
	p := truncate.DefaultPolicy()
	if p.MaxLen != 0 {
		t.Fatalf("expected MaxLen 0, got %d", p.MaxLen)
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("default policy should be valid: %v", err)
	}
}
