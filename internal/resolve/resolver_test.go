package resolve_test

import (
	"os"
	"testing"

	"github.com/yourorg/vaultpull/internal/resolve"
)

func TestResolvePlaceholderFromVars(t *testing.T) {
	r := resolve.New(map[string]string{"ENV": "production"})
	got, err := r.Resolve("secret/${ENV}/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/production/db" {
		t.Errorf("got %q, want %q", got, "secret/production/db")
	}
}

func TestResolveFallsBackToEnv(t *testing.T) {
	os.Setenv("VAULTPULL_REGION", "us-east-1")
	t.Cleanup(func() { os.Unsetenv("VAULTPULL_REGION") })

	r := resolve.New(nil)
	got, err := r.Resolve("secret/${VAULTPULL_REGION}/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/us-east-1/app" {
		t.Errorf("got %q, want %q", got, "secret/us-east-1/app")
	}
}

func TestResolveUnknownPlaceholderReturnsError(t *testing.T) {
	r := resolve.New(nil)
	_, err := r.Resolve("secret/${MISSING_VAR}/db")
	if err == nil {
		t.Fatal("expected error for unresolved placeholder, got nil")
	}
}

func TestResolveAllExpandsMultiplePaths(t *testing.T) {
	r := resolve.New(map[string]string{"APP": "api", "ENV": "staging"})
	paths := []string{"secret/${ENV}/${APP}", "secret/${ENV}/shared"}
	got, err := r.ResolveAll(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"secret/staging/api", "secret/staging/shared"}
	for i, g := range got {
		if g != want[i] {
			t.Errorf("index %d: got %q, want %q", i, g, want[i])
		}
	}
}

func TestResolveAllStopsOnFirstError(t *testing.T) {
	r := resolve.New(nil)
	paths := []string{"secret/static", "secret/${UNDEFINED}/path"}
	_, err := r.ResolveAll(paths)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolveNoPlaceholdersPassesThrough(t *testing.T) {
	r := resolve.New(nil)
	const plain = "secret/myapp/config"
	got, err := r.Resolve(plain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != plain {
		t.Errorf("got %q, want %q", got, plain)
	}
}
