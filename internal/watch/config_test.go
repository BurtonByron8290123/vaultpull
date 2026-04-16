package watch_test

import (
	"testing"
	"time"

	"github.com/yourorg/vaultpull/internal/watch"
)

func TestDefaultConfigIsValid(t *testing.T) {
	cfg := watch.DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
}

func TestValidateRejectsShortInterval(t *testing.T) {
	cfg := watch.DefaultConfig()
	cfg.Interval = 500 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for interval < 1s")
	}
}

func TestValidateRejectsEmptySnapshotPath(t *testing.T) {
	cfg := watch.DefaultConfig()
	cfg.SnapshotPath = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty snapshot_path")
	}
}

func TestDefaultInterval(t *testing.T) {
	cfg := watch.DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Fatalf("expected 30s, got %v", cfg.Interval)
	}
}
