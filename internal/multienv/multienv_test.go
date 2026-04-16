package multienv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/multienv"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "multienv-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestWriteAllCreatesFiles(t *testing.T) {
	dir := tempDir(t)
	targets := []multienv.Target{
		{Name: "dev", Path: ".env.dev"},
		{Name: "prod", Path: ".env.prod"},
	}
	w := multienv.New(dir, targets)
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := w.WriteAll(secrets); err != nil {
		t.Fatalf("WriteAll: %v", err)
	}
	for _, tgt := range targets {
		p := filepath.Join(dir, tgt.Path)
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected file %s to exist: %v", p, err)
		}
	}
}

func TestWriteAllFiltersKeys(t *testing.T) {
	dir := tempDir(t)
	targets := []multienv.Target{
		{Name: "dev", Path: ".env.dev", Keys: []string{"FOO"}},
	}
	w := multienv.New(dir, targets)
	secrets := map[string]string{"FOO": "bar", "SECRET": "hidden"}
	if err := w.WriteAll(secrets); err != nil {
		t.Fatalf("WriteAll: %v", err)
	}
	got, err := env.Parse(filepath.Join(dir, ".env.dev"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, ok := got["SECRET"]; ok {
		t.Error("expected SECRET to be filtered out")
	}
	if got["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", got["FOO"])
	}
}

func TestWriteAllEmptyKeysIncludesAll(t *testing.T) {
	dir := tempDir(t)
	targets := []multienv.Target{
		{Name: "all", Path: ".env.all"},
	}
	w := multienv.New(dir, targets)
	secrets := map[string]string{"A": "1", "B": "2"}
	if err := w.WriteAll(secrets); err != nil {
		t.Fatalf("WriteAll: %v", err)
	}
	got, err := env.Parse(filepath.Join(dir, ".env.all"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestWriteAllFilePermissions(t *testing.T) {
	dir := tempDir(t)
	targets := []multienv.Target{
		{Name: "sec", Path: ".env.sec"},
	}
	w := multienv.New(dir, targets)
	if err := w.WriteAll(map[string]string{"K": "v"}); err != nil {
		t.Fatalf("WriteAll: %v", err)
	}
	info, err := os.Stat(filepath.Join(dir, ".env.sec"))
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
