package template

import (
	"os"
	"testing"
)

func TestRenderBraceStyle(t *testing.T) {
	r := New(map[string]string{"ENV": "production", "REGION": "us-east-1"}, false)
	got, err := r.Render("secret/${ENV}/${REGION}/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "secret/production/us-east-1/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRenderDollarStyle(t *testing.T) {
	r := New(map[string]string{"APP": "myapp"}, false)
	got, err := r.Render("secret/$APP/creds")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/myapp/creds" {
		t.Errorf("got %q", got)
	}
}

func TestRenderUnresolvedReturnsError(t *testing.T) {
	r := New(map[string]string{}, false)
	_, err := r.Render("secret/${MISSING}/path")
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestRenderFallsBackToEnv(t *testing.T) {
	os.Setenv("VAULTPULL_TEST_VAR", "staging")
	defer os.Unsetenv("VAULTPULL_TEST_VAR")

	r := New(map[string]string{}, true)
	got, err := r.Render("secret/${VAULTPULL_TEST_VAR}/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/staging/db" {
		t.Errorf("got %q", got)
	}
}

func TestRenderVarsOverrideEnv(t *testing.T) {
	os.Setenv("VAULTPULL_PRIO", "from-env")
	defer os.Unsetenv("VAULTPULL_PRIO")

	r := New(map[string]string{"VAULTPULL_PRIO": "from-vars"}, true)
	got, err := r.Render("${VAULTPULL_PRIO}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "from-vars" {
		t.Errorf("vars should take precedence over env, got %q", got)
	}
}

func TestRenderAllReturnsAllPaths(t *testing.T) {
	r := New(map[string]string{"ENV": "dev"}, false)
	paths := []string{"secret/${ENV}/db", "secret/${ENV}/cache"}
	got, err := r.RenderAll(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "secret/dev/db" || got[1] != "secret/dev/cache" {
		t.Errorf("unexpected output: %v", got)
	}
}

func TestRenderAllStopsOnFirstError(t *testing.T) {
	r := New(map[string]string{}, false)
	_, err := r.RenderAll([]string{"ok/path", "bad/${MISSING}"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRenderNoVariablesUnchanged(t *testing.T) {
	r := New(map[string]string{}, false)
	const plain = "secret/static/path"
	got, err := r.Render(plain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != plain {
		t.Errorf("got %q, want %q", got, plain)
	}
}
