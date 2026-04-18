package transform_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/transform"
)

func TestApplyStripPrefix(t *testing.T) {
	tr, err := transform.New(transform.Policy{PrefixStrip: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := tr.Apply(map[string]string{"APP_DB_HOST": "localhost", "OTHER": "val"})
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected APP_ prefix to be stripped from APP_DB_HOST")
	}
	if _, ok := out["OTHER"]; !ok {
		t.Error("expected OTHER key to be preserved")
	}
}

func TestApplyKeyUpper(t *testing.T) {
	tr, _ := transform.New(transform.Policy{KeyCase: "upper"})
	out := tr.Apply(map[string]string{"db_host": "localhost"})
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected key to be uppercased")
	}
}

func TestApplyKeyLower(t *testing.T) {
	tr, _ := transform.New(transform.Policy{KeyCase: "lower"})
	out := tr.Apply(map[string]string{"DB_HOST": "localhost"})
	if _, ok := out["db_host"]; !ok {
		t.Error("expected key to be lowercased")
	}
}

func TestApplyValueTrimSpace(t *testing.T) {
	tr, _ := transform.New(transform.Policy{ValueTrimSpace: true})
	out := tr.Apply(map[string]string{"KEY": "  value  "})
	if out["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestApplyNoOpPolicy(t *testing.T) {
	tr, _ := transform.New(transform.Policy{})
	in := map[string]string{"Key": " val "}
	out := tr.Apply(in)
	if out["Key"] != " val " {
		t.Errorf("expected unchanged value, got %q", out["Key"])
	}
}

func TestInvalidKeyCaseReturnsError(t *testing.T) {
	_, err := transform.New(transform.Policy{KeyCase: "title"})
	if err == nil {
		t.Fatal("expected error for invalid key_case")
	}
}

func TestApplyCombined(t *testing.T) {
	tr, _ := transform.New(transform.Policy{
		PrefixStrip:    "SVC_",
		KeyCase:        "lower",
		ValueTrimSpace: true,
	})
	out := tr.Apply(map[string]string{"SVC_PORT": "  8080  "})
	if out["port"] != "8080" {
		t.Errorf("expected 'port'='8080', got %v", out)
	}
}
