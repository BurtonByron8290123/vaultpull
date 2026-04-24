package envresolve

import (
	"os"
	"testing"
)

func newResolver(allowFallback, errorOnMissing bool, overrides map[string]string) *Resolver {
	return New(Policy{AllowEnvFallback: allowFallback, ErrorOnMissing: errorOnMissing}, overrides)
}

func TestApplyExpandsBraceStyle(t *testing.T) {
	r := newResolver(false, false, nil)
	in := map[string]string{
		"BASE": "/opt/app",
		"BIN":  "${BASE}/bin",
	}
	out, err := r.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["BIN"]; got != "/opt/app/bin" {
		t.Errorf("BIN = %q, want %q", got, "/opt/app/bin")
	}
}

func TestApplyExpandsDollarStyle(t *testing.T) {
	r := newResolver(false, false, nil)
	in := map[string]string{
		"HOST": "localhost",
		"URL":  "http://$HOST:8080",
	}
	out, err := r.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["URL"]; got != "http://localhost:8080" {
		t.Errorf("URL = %q, want %q", got, "http://localhost:8080")
	}
}

func TestApplyOverridesTakePrecedence(t *testing.T) {
	r := newResolver(false, false, map[string]string{"HOST": "override-host"})
	in := map[string]string{
		"HOST": "original-host",
		"URL":  "http://${HOST}",
	}
	out, err := r.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["URL"]; got != "http://override-host" {
		t.Errorf("URL = %q, want %q", got, "http://override-host")
	}
}

func TestApplyFallsBackToProcessEnv(t *testing.T) {
	t.Setenv("PROC_VAR", "from-process")
	r := newResolver(true, false, nil)
	in := map[string]string{"VAL": "${PROC_VAR}"}
	out, err := r.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["VAL"]; got != "from-process" {
		t.Errorf("VAL = %q, want %q", got, "from-process")
	}
}

func TestApplyLeavesUnresolvedWhenNotErrorOnMissing(t *testing.T) {
	r := newResolver(false, false, nil)
	in := map[string]string{"VAL": "${MISSING_KEY}"}
	out, err := r.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["VAL"]; got != "${MISSING_KEY}" {
		t.Errorf("VAL = %q, want original placeholder", got)
	}
}

func TestApplyErrorOnMissing(t *testing.T) {
	os.Unsetenv("DEFINITELY_ABSENT")
	r := newResolver(false, true, nil)
	in := map[string]string{"VAL": "${DEFINITELY_ABSENT}"}
	_, err := r.Apply(in)
	if err == nil {
		t.Fatal("expected error for unresolved placeholder, got nil")
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	r := newResolver(false, false, nil)
	in := map[string]string{
		"A": "hello",
		"B": "${A} world",
	}
	orig := in["B"]
	_, _ = r.Apply(in)
	if in["B"] != orig {
		t.Error("Apply mutated the input map")
	}
}
