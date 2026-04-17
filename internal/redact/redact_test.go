package redact_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/redact"
)

func TestIsSensitiveMatchesCaseInsensitive(t *testing.T) {
	r := redact.New([]string{"API_KEY", "password"})
	if !r.IsSensitive("api_key") {
		t.Error("expected api_key to be sensitive")
	}
	if !r.IsSensitive("PASSWORD") {
		t.Error("expected PASSWORD to be sensitive")
	}
	if r.IsSensitive("USERNAME") {
		t.Error("expected USERNAME to not be sensitive")
	}
}

func TestValueMasksSensitiveKey(t *testing.T) {
	r := redact.New([]string{"SECRET"})
	got := r.Value("SECRET", "super-secret")
	if got != "****" {
		t.Errorf("expected masked value, got %q", got)
	}
}

func TestValuePassesThroughNonSensitive(t *testing.T) {
	r := redact.New([]string{"SECRET"})
	got := r.Value("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestWithMaskUsesCustomMask(t *testing.T) {
	r := redact.New([]string{"TOKEN"}).WithMask("[REDACTED]")
	got := r.Value("TOKEN", "abc123")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestMapRedactsOnlySensitiveKeys(t *testing.T) {
	r := redact.New([]string{"DB_PASSWORD"})
	input := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "s3cr3t",
		"DB_PORT":     "5432",
	}
	out := r.Map(input)
	if out["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST should be unchanged")
	}
	if out["DB_PASSWORD"] != "****" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT should be unchanged")
	}
}

func TestMapDoesNotMutateOriginal(t *testing.T) {
	r := redact.New([]string{"SECRET"})
	input := map[string]string{"SECRET": "value"}
	r.Map(input)
	if input["SECRET"] != "value" {
		t.Error("original map should not be mutated")
	}
}
