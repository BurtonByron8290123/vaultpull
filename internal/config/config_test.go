package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("VAULTPULL_VAULT_TOKEN", "test-token")

	cfg := mustLoadInline(t, `
paths:
  - secret/app/db
`)
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("default vault_addr wrong: %q", cfg.VaultAddr)
	}
	if cfg.MaxBackups != 5 {
		t.Errorf("default max_backups wrong: %d", cfg.MaxBackups)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("default output_file wrong: %q", cfg.OutputFile)
	}
	if !cfg.Merge {
		t.Error("default merge should be true")
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Setenv("VAULTPULL_VAULT_TOKEN", "")

	cfg := mustLoadInline(t, `
vault_token: file-token
vault_addr: https://vault.example.com
paths:
  - secret/prod/db
  - secret/prod/cache
output_file: prod.env
template_vars:
  ENV: prod
`)
	if cfg.VaultToken != "file-token" {
		t.Errorf("vault_token: got %q", cfg.VaultToken)
	}
	if cfg.VaultAddr != "https://vault.example.com" {
		t.Errorf("vault_addr: got %q", cfg.VaultAddr)
	}
	if len(cfg.Paths) != 2 {
		t.Errorf("paths len: got %d", len(cfg.Paths))
	}
	if cfg.TemplateVars["ENV"] != "prod" {
		t.Errorf("template_vars[ENV]: got %q", cfg.TemplateVars["ENV"])
	}
}

func TestValidateMissingToken(t *testing.T) {
	os.Unsetenv("VAULTPULL_VAULT_TOKEN")
	tmp := writeConfigFile(t, `
paths:
  - secret/app
`)
	_, err := Load(tmp)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestValidateMissingPath(t *testing.T) {
	t.Setenv("VAULTPULL_VAULT_TOKEN", "tok")
	tmp := writeConfigFile(t, `vault_token: tok\n`)
	_, err := Load(tmp)
	if err == nil {
		t.Fatal("expected error for missing paths")
	}
}

// helpers

func mustLoadInline(t *testing.T, yaml string) *Config {
	t.Helper()
	tmp := writeConfigFile(t, yaml)
	cfg, err := Load(tmp)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return cfg
}

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".vaultpull.yaml")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}
