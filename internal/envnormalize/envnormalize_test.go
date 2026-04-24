package envnormalize_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envnormalize"
)

func newNormalizer(t *testing.T, p envnormalize.Policy) *envnormalize.Normalizer {
	t.Helper()
	n, err := envnormalize.New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return n
}

func TestUpperCaseKeys(t *testing.T) {
	p := envnormalize.DefaultPolicy()
	n := newNormalizer(t, p)
	out := n.Apply(map[string]string{"db_host": "localhost", "api_key": "secret"})
	for _, k := range []string{"DB_HOST", "API_KEY"} {
		if _, ok := out[k]; !ok {
			t.Errorf("expected key %q in output", k)
		}
	}
}

func TestLowerCaseKeys(t *testing.T) {
	p := envnormalize.Policy{KeyStrategy: envnormalize.StrategyLower, TrimValues: false}
	n := newNormalizer(t, p)
	out := n.Apply(map[string]string{"DB_HOST": "localhost"})
	if _, ok := out["db_host"]; !ok {
		t.Error("expected lower-cased key db_host")
	}
}

func TestNoneStrategyPreservesKeys(t *testing.T) {
	p := envnormalize.Policy{KeyStrategy: envnormalize.StrategyNone}
	n := newNormalizer(t, p)
	out := n.Apply(map[string]string{"MixedCase": "val"})
	if _, ok := out["MixedCase"]; !ok {
		t.Error("expected original key MixedCase to be preserved")
	}
}

func TestTrimValuesRemovesWhitespace(t *testing.T) {
	p := envnormalize.DefaultPolicy()
	n := newNormalizer(t, p)
	out := n.Apply(map[string]string{"KEY": "  value  "})
	if got := out["KEY"]; got != "value" {
		t.Errorf("expected trimmed value %q, got %q", "value", got)
	}
}

func TestTrimDisabledPreservesWhitespace(t *testing.T) {
	p := envnormalize.Policy{KeyStrategy: envnormalize.StrategyNone, TrimValues: false}
	n := newNormalizer(t, p)
	out := n.Apply(map[string]string{"KEY": "  value  "})
	if got := out["KEY"]; got != "  value  " {
		t.Errorf("expected untrimmed value, got %q", got)
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	p := envnormalize.DefaultPolicy()
	n := newNormalizer(t, p)
	in := map[string]string{"lower_key": "  val  "}
	_ = n.Apply(in)
	if _, ok := in["lower_key"]; !ok {
		t.Error("Apply mutated the input map")
	}
	if in["lower_key"] != "  val  " {
		t.Error("Apply mutated the value in the input map")
	}
}

func TestApplyEmptyMapReturnsEmptyMap(t *testing.T) {
	p := envnormalize.DefaultPolicy()
	n := newNormalizer(t, p)
	out := n.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %v", out)
	}
}
