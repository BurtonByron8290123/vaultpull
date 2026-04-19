package envexport_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envexport"
)

func TestNewRejectsUnknownFormat(t *testing.T) {
	_, err := envexport.New("xml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestWriteDotenv(t *testing.T) {
	ex, err := envexport.New(envexport.FormatDotenv)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	var buf bytes.Buffer
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"}
	if err := ex.Write(&buf, secrets); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_PASS=secret") {
		t.Errorf("missing DB_PASS line, got:\n%s", out)
	}
	if strings.Contains(out, "export ") {
		t.Errorf("dotenv format should not include 'export', got:\n%s", out)
	}
}

func TestWriteExportFormat(t *testing.T) {
	ex, _ := envexport.New(envexport.FormatExport)
	var buf bytes.Buffer
	_ = ex.Write(&buf, map[string]string{"TOKEN": "xyz"})
	if !strings.HasPrefix(buf.String(), "export ") {
		t.Errorf("expected 'export ' prefix, got: %s", buf.String())
	}
}

func TestWriteJSON(t *testing.T) {
	ex, _ := envexport.New(envexport.FormatJSON)
	var buf bytes.Buffer
	secrets := map[string]string{"FOO": "bar"}
	if err := ex.Write(&buf, secrets); err != nil {
		t.Fatalf("Write: %v", err)
	}
	var got map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %v", got)
	}
}

func TestFromEnvUsesDefault(t *testing.T) {
	t.Setenv("VAULTPULL_EXPORT_FORMAT", "")
	cfg, err := envexport.FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if cfg.Format != envexport.FormatDotenv {
		t.Errorf("expected dotenv default, got %q", cfg.Format)
	}
}

func TestFromEnvReadsFormat(t *testing.T) {
	t.Setenv("VAULTPULL_EXPORT_FORMAT", "json")
	cfg, err := envexport.FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if cfg.Format != envexport.FormatJSON {
		t.Errorf("expected json, got %q", cfg.Format)
	}
}

func TestFromEnvInvalidFormatReturnsError(t *testing.T) {
	t.Setenv("VAULTPULL_EXPORT_FORMAT", "toml")
	_, err := envexport.FromEnv()
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}
