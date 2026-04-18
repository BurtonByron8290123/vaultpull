package mask_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/mask"
)

func TestApplyFullMask(t *testing.T) {
	p := mask.DefaultPolicy()
	got := p.Apply("supersecret")
	if got != "********" {
		t.Fatalf("expected masked value, got %q", got)
	}
}

func TestApplyEmptyValueUnchanged(t *testing.T) {
	p := mask.DefaultPolicy()
	if got := p.Apply(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestApplyRevealTrailingChars(t *testing.T) {
	p := mask.Policy{Mask: "****", RevealChars: 3}
	got := p.Apply("abcdef")
	if got != "****def" {
		t.Fatalf("expected ****def, got %q", got)
	}
}

func TestApplyRevealCharsExceedsLength(t *testing.T) {
	p := mask.Policy{Mask: "****", RevealChars: 20}
	got := p.Apply("short")
	if got != "****" {
		t.Fatalf("expected full mask when reveal >= len, got %q", got)
	}
}

func TestApplyMapMasksAllValues(t *testing.T) {
	p := mask.DefaultPolicy()
	secrets := map[string]string{"A": "val1", "B": "val2"}
	out := p.ApplyMap(secrets)
	for k, v := range out {
		if v != "********" {
			t.Errorf("key %s: expected masked, got %q", k, v)
		}
	}
}

func TestApplyMapDoesNotMutateOriginal(t *testing.T) {
	p := mask.DefaultPolicy()
	orig := map[string]string{"KEY": "plaintext"}
	p.ApplyMap(orig)
	if orig["KEY"] != "plaintext" {
		t.Fatal("original map was mutated")
	}
}

func TestIsSensitiveMatchesKnownPatterns(t *testing.T) {
	cases := []string{"DB_PASSWORD", "api_key", "AUTH_TOKEN", "PRIVATE_KEY", "aws_secret"}
	for _, c := range cases {
		if !mask.IsSensitive(c) {
			t.Errorf("expected %q to be sensitive", c)
		}
	}
}

func TestIsSensitiveIgnoresNonSensitiveKeys(t *testing.T) {
	cases := []string{"HOST", "PORT", "LOG_LEVEL", "APP_NAME"}
	for _, c := range cases {
		if mask.IsSensitive(c) {
			t.Errorf("expected %q to NOT be sensitive", c)
		}
	}
}
