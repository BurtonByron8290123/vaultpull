package prefixmap_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/prefixmap"
)

func TestApplyMatchingEntry(t *testing.T) {
	m, _ := prefixmap.New([]prefixmap.Entry{{Path: "secret/app", Prefix: "APP_"}})
	got := m.Apply("secret/app", "DB_HOST")
	if got != "APP_DB_HOST" {
		t.Fatalf("expected APP_DB_HOST, got %s", got)
	}
}

func TestApplyNoMatch(t *testing.T) {
	m, _ := prefixmap.New([]prefixmap.Entry{{Path: "secret/app", Prefix: "APP_"}})
	got := m.Apply("secret/other", "KEY")
	if got != "KEY" {
		t.Fatalf("expected KEY unchanged, got %s", got)
	}
}

func TestApplyEmptyPrefix(t *testing.T) {
	m, _ := prefixmap.New([]prefixmap.Entry{{Path: "secret/app", Prefix: ""}})
	got := m.Apply("secret/app", "KEY")
	if got != "KEY" {
		t.Fatalf("expected KEY unchanged, got %s", got)
	}
}

func TestApplyMapRewrites(t *testing.T) {
	m, _ := prefixmap.New([]prefixmap.Entry{{Path: "secret/svc", Prefix: "SVC_"}})
	out := m.ApplyMap("secret/svc", map[string]string{"HOST": "localhost", "PORT": "5432"})
	if out["SVC_HOST"] != "localhost" || out["SVC_PORT"] != "5432" {
		t.Fatalf("unexpected map: %v", out)
	}
}

func TestNewRejectsEmptyPath(t *testing.T) {
	_, err := prefixmap.New([]prefixmap.Entry{{Path: "", Prefix: "X_"}})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "map.yaml")
	content := "mappings:\n  - path: secret/app\n    prefix: APP_\n"
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	m, err := prefixmap.LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Apply("secret/app", "TOKEN"); got != "APP_TOKEN" {
		t.Fatalf("expected APP_TOKEN, got %s", got)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := prefixmap.LoadConfig("/nonexistent/map.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
