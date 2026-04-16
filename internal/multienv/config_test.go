package multienv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/multienv"
)

func writeYAML(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return p
}

func TestLoadConfigValid(t *testing.T) {
	dir := tempDir(t)
	p := writeYAML(t, dir, "multienv.yaml", `
targets:
  - name: dev
    path: .env.dev
    keys: [FOO, BAR]
  - name: prod
    path: .env.prod
`)
	cfg, err := multienv.LoadConfig(p)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if len(cfg.Targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(cfg.Targets))
	}
	if cfg.Targets[0].Name != "dev" {
		t.Errorf("expected name=dev, got %q", cfg.Targets[0].Name)
	}
	if len(cfg.Targets[0].Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(cfg.Targets[0].Keys))
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := multienv.LoadConfig("/nonexistent/multienv.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfigMissingName(t *testing.T) {
	dir := tempDir(t)
	p := writeYAML(t, dir, "bad.yaml", `
targets:
  - path: .env.dev
`)
	_, err := multienv.LoadConfig(p)
	if err == nil {
		t.Fatal("expected validation error for missing name")
	}
}

func TestLoadConfigMissingPath(t *testing.T) {
	dir := tempDir(t)
	p := writeYAML(t, dir, "bad.yaml", `
targets:
  - name: dev
`)
	_, err := multienv.LoadConfig(p)
	if err == nil {
		t.Fatal("expected validation error for missing path")
	}
}
