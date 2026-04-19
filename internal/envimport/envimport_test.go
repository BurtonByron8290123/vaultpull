package envimport_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultpull/internal/envimport"
)

func writeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeEnv: %v", err)
	}
	return p
}

func TestImportAddsNewKeys(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "FOO=bar\nBAZ=qux\n")

	im := envimport.New(envimport.DefaultPolicy())
	got, err := im.Import(src, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestImportSkipExistingByDefault(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "FOO=from_file\n")

	base := map[string]string{"FOO": "from_base"}
	im := envimport.New(envimport.DefaultPolicy())
	got, err := im.Import(src, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "from_base" {
		t.Errorf("expected base value to be preserved, got %q", got["FOO"])
	}
}

func TestImportOverwritesWhenSkipExistingFalse(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "FOO=from_file\n")

	base := map[string]string{"FOO": "from_base"}
	im := envimport.New(envimport.Policy{SkipExisting: false})
	got, err := im.Import(src, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "from_file" {
		t.Errorf("expected file value, got %q", got["FOO"])
	}
}

func TestImportMissingFileReturnsError(t *testing.T) {
	im := envimport.New(envimport.DefaultPolicy())
	_, err := im.Import("/nonexistent/.env", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestImportIgnoreMissingFile(t *testing.T) {
	base := map[string]string{"KEY": "val"}
	im := envimport.New(envimport.Policy{IgnoreMissing: true, SkipExisting: true})
	got, err := im.Import("/nonexistent/.env", base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "val" {
		t.Errorf("expected base to be returned unchanged")
	}
}

func TestImportDoesNotMutateBase(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "NEW=value\n")

	base := map[string]string{"EXISTING": "yes"}
	im := envimport.New(envimport.DefaultPolicy())
	_, err := im.Import(src, base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := base["NEW"]; ok {
		t.Error("base map was mutated")
	}
}
