package envclone_test

import (
	"os"
	"testing"

	"github.com/example/vaultpull/internal/envclone"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_CLONE_OVERWRITE")
	os.Unsetenv("VAULTPULL_CLONE_DRY_RUN")
	p := envclone.FromEnv()
	if p.Overwrite {
		t.Error("expected Overwrite=false by default")
	}
	if p.DryRun {
		t.Error("expected DryRun=false by default")
	}
}

func TestFromEnvReadsOverwrite(t *testing.T) {
	t.Setenv("VAULTPULL_CLONE_OVERWRITE", "true")
	p := envclone.FromEnv()
	if !p.Overwrite {
		t.Error("expected Overwrite=true")
	}
}

func TestFromEnvReadsDryRun(t *testing.T) {
	t.Setenv("VAULTPULL_CLONE_DRY_RUN", "1")
	p := envclone.FromEnv()
	if !p.DryRun {
		t.Error("expected DryRun=true")
	}
}
