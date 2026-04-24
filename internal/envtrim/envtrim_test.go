package envtrim_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envtrim"
)

func newTrimmer(t *testing.T, p envtrim.Policy) *envtrim.Trimmer {
	t.Helper()
	tr, err := envtrim.New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestApplyRemovesBlankValues(t *testing.T) {
	tr := newTrimmer(t, envtrim.Policy{StripBlank: true})
	env := map[string]string{"A": "hello", "B": "", "C": "   "}
	out := tr.Apply(env)
	if _, ok := out["B"]; ok {
		t.Error("expected key B to be removed")
	}
	if _, ok := out["C"]; ok {
		t.Error("expected key C (whitespace) to be removed")
	}
	if out["A"] != "hello" {
		t.Errorf("expected A=hello, got %q", out["A"])
	}
}

func TestApplyKeepsBlankValuesWhenDisabled(t *testing.T) {
	tr := newTrimmer(t, envtrim.Policy{StripBlank: false})
	env := map[string]string{"A": ""}
	out := tr.Apply(env)
	if _, ok := out["A"]; !ok {
		t.Error("expected key A to be kept when StripBlank is false")
	}
}

func TestApplyRemovesSentinelValues(t *testing.T) {
	tr := newTrimmer(t, envtrim.Policy{StripSentinels: true})
	env := map[string]string{
		"A": "null",
		"B": "NIL",
		"C": "None",
		"D": "real-value",
	}
	out := tr.Apply(env)
	for _, k := range []string{"A", "B", "C"} {
		if _, ok := out[k]; ok {
			t.Errorf("expected key %s to be removed as sentinel", k)
		}
	}
	if out["D"] != "real-value" {
		t.Errorf("expected D=real-value, got %q", out["D"])
	}
}

func TestApplyExtraSentinels(t *testing.T) {
	tr := newTrimmer(t, envtrim.Policy{
		StripSentinels: true,
		Extra:          []string{"PLACEHOLDER", "tbd"},
	})
	env := map[string]string{"X": "placeholder", "Y": "TBD", "Z": "keep"}
	out := tr.Apply(env)
	for _, k := range []string{"X", "Y"} {
		if _, ok := out[k]; ok {
			t.Errorf("expected key %s to be trimmed via Extra sentinels", k)
		}
	}
	if out["Z"] != "keep" {
		t.Errorf("expected Z=keep, got %q", out["Z"])
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	tr := newTrimmer(t, envtrim.DefaultPolicy())
	env := map[string]string{"A": "", "B": "val"}
	_ = tr.Apply(env)
	if _, ok := env["A"]; !ok {
		t.Error("Apply must not mutate the input map")
	}
}

func TestDefaultPolicyStripsBlankAndSentinels(t *testing.T) {
	p := envtrim.DefaultPolicy()
	if !p.StripBlank {
		t.Error("DefaultPolicy should have StripBlank=true")
	}
	if !p.StripSentinels {
		t.Error("DefaultPolicy should have StripSentinels=true")
	}
}
