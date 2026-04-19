package labelmap_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/labelmap"
)

func TestFromEnvEmptyVarReturnsNoop(t *testing.T) {
	os.Unsetenv("VAULTPULL_LABELMAP_FILE")
	m, err := labelmap.FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	out := m.Apply(map[string]string{"K": "V"})
	if out["K"] != "V" {
		t.Error("noop mapper should preserve keys")
	}
}

func TestFromEnvLoadsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lm.yaml")
	os.WriteFile(path, []byte("- from: OLD\n  to: NEW\n"), 0o600)
	t.Setenv("VAULTPULL_LABELMAP_FILE", path)

	m, err := labelmap.FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	out := m.Apply(map[string]string{"OLD": "val"})
	if out["NEW"] != "val" {
		t.Errorf("expected NEW=val, got %v", out)
	}
}

func TestFromEnvInvalidFileReturnsError(t *testing.T) {
	t.Setenv("VAULTPULL_LABELMAP_FILE", "/no/such/file.yaml")
	_, err := labelmap.FromEnv()
	if err == nil {
		t.Error("expected error for missing file")
	}
}
