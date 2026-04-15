package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "test-token")
	t.Setenv("VAULTPULL_VAULT_PATH", "secret/myapp")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.OutputFile != ".env" {
		t.Errorf("expected default output_file '.env', got %q", cfg.OutputFile)
	}
	if cfg.BackupDir != ".env.backups" {
		t.Errorf("expected default backup_dir '.env.backups', got %q", cfg.BackupDir)
	}
	if cfg.Rotate != false {
		t.Errorf("expected default rotate false, got %v", cfg.Rotate)
	}
	if cfg.VaultToken != "test-token" {
		t.Errorf("expected vault_token 'test-token', got %q", cfg.VaultToken)
	}
}

func TestLoadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfgContent := []byte(`
vault_addr: "http://vault.example.com:8200"
vault_token: "s.mytoken"
vault_path: "secret/data/myapp"
output_file: ".env.production"
rotate: true
backup_dir: ".backups"
mappings:
  DB_PASSWORD: database_password
  API_KEY: api_key
`)
	cfgFile := filepath.Join(tmpDir, ".vaultpull.yaml")
	if err := os.WriteFile(cfgFile, cfgContent, 0600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := Load(cfgFile)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.VaultAddr != "http://vault.example.com:8200" {
		t.Errorf("unexpected vault_addr: %q", cfg.VaultAddr)
	}
	if cfg.VaultPath != "secret/data/myapp" {
		t.Errorf("unexpected vault_path: %q", cfg.VaultPath)
	}
	if cfg.OutputFile != ".env.production" {
		t.Errorf("unexpected output_file: %q", cfg.OutputFile)
	}
	if !cfg.Rotate {
		t.Error("expected rotate to be true")
	}
	if len(cfg.Mappings) != 2 {
		t.Errorf("expected 2 mappings, got %d", len(cfg.Mappings))
	}
}

func TestValidateMissingToken(t *testing.T) {
	cfg := &Config{
		VaultAddr: "http://127.0.0.1:8200",
		VaultPath: "secret/myapp",
	}
	if err := cfg.validate(); err == nil {
		t.Error("expected validation error for missing token")
	}
}

func TestValidateMissingPath(t *testing.T) {
	cfg := &Config{
		VaultAddr:  "http://127.0.0.1:8200",
		VaultToken: "s.token",
	}
	if err := cfg.validate(); err == nil {
		t.Error("expected validation error for missing vault_path")
	}
}
