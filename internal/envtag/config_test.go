package envtag_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpull/internal/envtag"
)

func writeJSON(t *testing.T, v any) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	path := filepath.Join(t.TempDir(), "rules.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}

func TestFromFileParsesRules(t *testing.T) {
	raw := []map[string]any{
		{"prefix": "DB_", "tags": map[string]string{"group": "database"}},
	}
	path := writeJSON(t, raw)
	p, err := envtag.FromFile(path)
	if err != nil {
		t.Fatalf("FromFile: %v", err)
	}
	if len(p.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(p.Rules))
	}
	if p.Rules[0].Prefix != "DB_" {
		t.Errorf("unexpected prefix: %q", p.Rules[0].Prefix)
	}
	if len(p.Rules[0].Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(p.Rules[0].Tags))
	}
}

func TestFromFileMissingReturnsError(t *testing.T) {
	_, err := envtag.FromFile("/nonexistent/rules.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFromFileBadJSONReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o600)
	_, err := envtag.FromFile(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestFromEnvEmptyVarReturnsEmptyPolicy(t *testing.T) {
	t.Setenv("VAULTPULL_ENVTAG_CONFIG", "")
	p, err := envtag.FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if len(p.Rules) != 0 {
		t.Errorf("expected empty policy")
	}
}

func TestFromEnvLoadsFile(t *testing.T) {
	raw := []map[string]any{
		{"prefix": "SVC_", "tags": map[string]string{"tier": "backend"}},
	}
	path := writeJSON(t, raw)
	t.Setenv("VAULTPULL_ENVTAG_CONFIG", path)
	p, err := envtag.FromEnv()
	if err != nil {
		t.Fatalf("FromEnv: %v", err)
	}
	if len(p.Rules) != 1 || p.Rules[0].Prefix != "SVC_" {
		t.Errorf("unexpected policy: %+v", p)
	}
}
