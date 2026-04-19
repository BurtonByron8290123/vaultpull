package scope_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/scope"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv("VAULTPULL_SCOPE_ALLOW", "")
	t.Setenv("VAULTPULL_SCOPE_DENY", "")
	p := scope.FromEnv()
	if len(p.Allow) != 0 || len(p.Deny) != 0 {
		t.Errorf("expected empty policy, got %+v", p)
	}
}

func TestFromEnvReadsAllow(t *testing.T) {
	t.Setenv("VAULTPULL_SCOPE_ALLOW", "secret/app, secret/shared")
	t.Setenv("VAULTPULL_SCOPE_DENY", "")
	p := scope.FromEnv()
	if len(p.Allow) != 2 {
		t.Fatalf("expected 2 allow entries, got %d", len(p.Allow))
	}
	if p.Allow[0] != "secret/app" || p.Allow[1] != "secret/shared" {
		t.Errorf("unexpected allow values: %v", p.Allow)
	}
}

func TestFromEnvReadsDeny(t *testing.T) {
	t.Setenv("VAULTPULL_SCOPE_ALLOW", "")
	t.Setenv("VAULTPULL_SCOPE_DENY", "secret/internal")
	p := scope.FromEnv()
	if len(p.Deny) != 1 || p.Deny[0] != "secret/internal" {
		t.Errorf("unexpected deny values: %v", p.Deny)
	}
}

func TestFromEnvSkipsBlanks(t *testing.T) {
	t.Setenv("VAULTPULL_SCOPE_ALLOW", "secret/app,,secret/shared,")
	t.Setenv("VAULTPULL_SCOPE_DENY", "")
	p := scope.FromEnv()
	if len(p.Allow) != 2 {
		t.Fatalf("expected 2 entries after blank skip, got %d", len(p.Allow))
	}
}
