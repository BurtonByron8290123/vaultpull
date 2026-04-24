package envtemplate

import (
	"testing"
)

func newRenderer(t *testing.T, p Policy) *Renderer {
	t.Helper()
	r, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return r
}

func TestApplyRendersSimpleReference(t *testing.T) {
	r := newRenderer(t, DefaultPolicy())
	src := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://{{index . \"HOST\"}}/db",
	}
	out, err := r.Apply(src)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got, want := out["DSN"], "postgres://localhost/db"; got != want {
		t.Errorf("DSN = %q, want %q", got, want)
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	r := newRenderer(t, DefaultPolicy())
	src := map[string]string{"A": "hello", "B": "world"}
	_, err := r.Apply(src)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if src["A"] != "hello" || src["B"] != "world" {
		t.Error("Apply mutated the input map")
	}
}

func TestApplyMissingKeyZeroByDefault(t *testing.T) {
	r := newRenderer(t, DefaultPolicy())
	src := map[string]string{
		"VAL": "prefix-{{index . \"MISSING\"}}-suffix",
	}
	out, err := r.Apply(src)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got, want := out["VAL"], "prefix--suffix"; got != want {
		t.Errorf("VAL = %q, want %q", got, want)
	}
}

func TestApplyErrorOnMissingKey(t *testing.T) {
	p := DefaultPolicy()
	p.ErrorOnMissing = true
	r := newRenderer(t, p)
	src := map[string]string{
		"VAL": "{{index . \"GHOST\"}}",
	}
	_, err := r.Apply(src)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestApplyCustomDelimiters(t *testing.T) {
	p := Policy{LeftDelim: "((", RightDelim: "))", ErrorOnMissing: false}
	r := newRenderer(t, p)
	src := map[string]string{
		"NAME": "world",
		"MSG":  "hello ((index . \"NAME\"))",
	}
	out, err := r.Apply(src)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got, want := out["MSG"], "hello world"; got != want {
		t.Errorf("MSG = %q, want %q", got, want)
	}
}

func TestNewRejectsEmptyDelimiters(t *testing.T) {
	p := Policy{LeftDelim: "", RightDelim: "}}"}
	_, err := New(p)
	if err == nil {
		t.Fatal("expected error for empty left delimiter")
	}
}

func TestApplyPlainValuePassesThrough(t *testing.T) {
	r := newRenderer(t, DefaultPolicy())
	src := map[string]string{"PLAIN": "no-template-here"}
	out, err := r.Apply(src)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got, want := out["PLAIN"], "no-template-here"; got != want {
		t.Errorf("PLAIN = %q, want %q", got, want)
	}
}
