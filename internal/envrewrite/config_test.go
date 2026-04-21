package envrewrite

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigMissingFileReturnsError(t *testing.T) {
	_, err := LoadConfig("/no/such/file.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadConfigInvalidYAMLReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte(":::invalid"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadConfig(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadConfigMultipleRules(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.yaml")
	content := `rules:
  - pattern: "^dev"
    replacement: "prod"
    key_glob: "ENV_"
  - pattern: "8080"
    replacement: "443"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	p, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(p.Rules))
	}
	if p.Rules[0].KeyGlob != "ENV_" {
		t.Errorf("unexpected key_glob: %s", p.Rules[0].KeyGlob)
	}
}

func TestFromEnvLoadsFileFromEnvVar(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rw.yaml")
	content := "rules:\n  - pattern: \"old\"\n    replacement: \"new\"\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("VAULTPULL_REWRITE_CONFIG", path)
	rw, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	out := rw.Apply(map[string]string{"K": "old_value"})
	if out["K"] != "new_value" {
		t.Errorf("expected new_value, got %s", out["K"])
	}
}

func TestFromEnvInvalidFileReturnsError(t *testing.T) {
	t.Setenv("VAULTPULL_REWRITE_CONFIG", "/no/such/path.yaml")
	_, err := FromEnv()
	if err == nil {
		t.Error("expected error for missing config file")
	}
}
