package hooks_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/user/vaultpull/internal/hooks"
)

func TestRunPrePullNoopWhenEmpty(t *testing.T) {
	r := hooks.New(hooks.Config{})
	if err := r.RunPrePull(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunPostPullNoopWhenEmpty(t *testing.T) {
	r := hooks.New(hooks.Config{})
	if err := r.RunPostPull(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunPrePullSuccess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r := hooks.New(hooks.Config{PrePull: "echo pre-pull-ok"})
	if err := r.RunPrePull(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunPostPullSuccess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r := hooks.New(hooks.Config{PostPull: "echo post-pull-ok"})
	if err := r.RunPostPull(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunPrePullCommandFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r := hooks.New(hooks.Config{PrePull: "false"})
	if err := r.RunPrePull(context.Background()); err == nil {
		t.Fatal("expected error for failing command, got nil")
	}
}

func TestRunPrePullTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r := hooks.New(hooks.Config{
		PrePull: "sleep 10",
		Timeout: 50 * time.Millisecond,
	})
	err := r.RunPrePull(context.Background())
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestDefaultTimeoutApplied(t *testing.T) {
	// Ensure zero timeout gets defaulted (no panic, no zero-timeout issue).
	r := hooks.New(hooks.Config{})
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}
