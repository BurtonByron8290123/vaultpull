package preflight_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/preflight"
)

func TestRunNoChecksReturnsNil(t *testing.T) {
	r := preflight.New()
	if err := r.Run(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRunReturnsErrorOnFailure(t *testing.T) {
	fail := preflight.Check{
		Name: "always_fail",
		Fn:   func() error { return os.ErrPermission },
	}
	r := preflight.New(fail)
	if err := r.Run(); err == nil {
		t.Fatal("expected error")
	}
}

func TestRunCombinesMultipleErrors(t *testing.T) {
	mkFail := func(name string) preflight.Check {
		return preflight.Check{Name: name, Fn: func() error { return os.ErrNotExist }}
	}
	r := preflight.New(mkFail("a"), mkFail("b"))
	err := r.Run()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "a:") || !strings.Contains(err.Error(), "b:") {
		t.Errorf("expected both check names in error, got: %v", err)
	}
}

func TestCheckVaultAddrEmptyFails(t *testing.T) {
	c := preflight.CheckVaultAddr("")
	if err := c.Fn(); err == nil {
		t.Fatal("expected error for empty addr")
	}
}

func TestCheckVaultAddrSetPasses(t *testing.T) {
	c := preflight.CheckVaultAddr("http://127.0.0.1:8200")
	if err := c.Fn(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckTokenEmptyFails(t *testing.T) {
	c := preflight.CheckToken("")
	if err := c.Fn(); err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestCheckOutputWritableValidDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	c := preflight.CheckOutputWritable(path)
	if err := c.Fn(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckOutputWritableEmptyPathFails(t *testing.T) {
	c := preflight.CheckOutputWritable("")
	if err := c.Fn(); err == nil {
		t.Fatal("expected error for empty path")
	}
}
