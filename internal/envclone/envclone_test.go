package envclone_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/vaultpull/internal/env"
	"github.com/example/vaultpull/internal/envclone"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "envclone-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestCloneCreatesNewFile(t *testing.T) {
	dst := filepath.Join(tempDir(t), ".env")
	c, _ := envclone.New(envclone.DefaultPolicy())
	n, err := c.Clone(dst, map[string]string{"FOO": "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("want 1 written, got %d", n)
	}
	entries, _ := env.Parse(dst)
	m := env.ToMap(entries)
	if m["FOO"] != "bar" {
		t.Fatalf("want bar, got %s", m["FOO"])
	}
}

func TestCloneSkipsExistingByDefault(t *testing.T) {
	dst := filepath.Join(tempDir(t), ".env")
	env.WriteFile(dst, map[string]string{"FOO": "original"}, 0600)
	c, _ := envclone.New(envclone.DefaultPolicy())
	n, err := c.Clone(dst, map[string]string{"FOO": "new"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("want 0 written, got %d", n)
	}
	entries, _ := env.Parse(dst)
	m := env.ToMap(entries)
	if m["FOO"] != "original" {
		t.Fatalf("want original, got %s", m["FOO"])
	}
}

func TestCloneOverwriteReplacesExisting(t *testing.T) {
	dst := filepath.Join(tempDir(t), ".env")
	env.WriteFile(dst, map[string]string{"FOO": "original"}, 0600)
	p := envclone.DefaultPolicy()
	p.Overwrite = true
	c, _ := envclone.New(p)
	n, err := c.Clone(dst, map[string]string{"FOO": "new"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("want 1, got %d", n)
	}
	entries, _ := env.Parse(dst)
	m := env.ToMap(entries)
	if m["FOO"] != "new" {
		t.Fatalf("want new, got %s", m["FOO"])
	}
}

func TestCloneDryRunDoesNotWrite(t *testing.T) {
	dst := filepath.Join(tempDir(t), ".env")
	p := envclone.DefaultPolicy()
	p.DryRun = true
	c, _ := envclone.New(p)
	n, err := c.Clone(dst, map[string]string{"FOO": "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("want 1, got %d", n)
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		t.Fatal("expected file not to exist in dry-run mode")
	}
}

func TestNewRejectsZeroFileMode(t *testing.T) {
	p := envclone.DefaultPolicy()
	p.FileMode = 0
	_, err := envclone.New(p)
	if err == nil {
		t.Fatal("expected error for zero file mode")
	}
}
