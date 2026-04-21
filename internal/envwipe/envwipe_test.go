package envwipe_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultpull/internal/envwipe"
)

func tempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestNewRejectsEmptyKeys(t *testing.T) {
	_, err := envwipe.New(envwipe.Policy{})
	if err == nil {
		t.Fatal("expected error for empty keys, got nil")
	}
}

func TestApplyRemovesKey(t *testing.T) {
	w, _ := envwipe.New(envwipe.Policy{Keys: []string{"SECRET"}})
	input := map[string]string{"SECRET": "abc", "KEEP": "yes"}
	out, res := w.Apply(input)

	if _, ok := out["SECRET"]; ok {
		t.Error("expected SECRET to be removed")
	}
	if out["KEEP"] != "yes" {
		t.Error("expected KEEP to be preserved")
	}
	if len(res.Removed) != 1 || res.Removed[0] != "SECRET" {
		t.Errorf("unexpected Removed: %v", res.Removed)
	}
}

func TestApplySkipsMissingKey(t *testing.T) {
	w, _ := envwipe.New(envwipe.Policy{Keys: []string{"MISSING"}})
	input := map[string]string{"OTHER": "val"}
	_, res := w.Apply(input)

	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("unexpected Skipped: %v", res.Skipped)
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	w, _ := envwipe.New(envwipe.Policy{Keys: []string{"A"}})
	input := map[string]string{"A": "1", "B": "2"}
	w.Apply(input)
	if _, ok := input["A"]; !ok {
		t.Error("Apply must not mutate the input map")
	}
}

func TestDryRunDoesNotModifyFile(t *testing.T) {
	path := tempEnvFile(t, "SECRET=abc\nKEEP=yes\n")
	w, _ := envwipe.New(envwipe.Policy{Keys: []string{"SECRET"}, DryRun: true})
	res, err := w.WipeFile(path)
	if err != nil {
		t.Fatalf("WipeFile: %v", err)
	}
	if len(res.Removed) != 1 {
		t.Errorf("expected 1 removed in dry-run result, got %d", len(res.Removed))
	}
	data, _ := os.ReadFile(path)
	if string(data) != "SECRET=abc\nKEEP=yes\n" {
		t.Error("dry-run must not modify the file")
	}
}

func TestWipeFileRemovesKey(t *testing.T) {
	path := tempEnvFile(t, "SECRET=abc\nKEEP=yes\n")
	w, _ := envwipe.New(envwipe.Policy{Keys: []string{"SECRET"}})
	_, err := w.WipeFile(path)
	if err != nil {
		t.Fatalf("WipeFile: %v", err)
	}
	data, _ := os.ReadFile(path)
	if contains(string(data), "SECRET") {
		t.Error("expected SECRET to be absent from file after wipe")
	}
}

func TestSummaryFormat(t *testing.T) {
	res := envwipe.Result{Removed: []string{"A", "B"}, Skipped: []string{"C"}}
	got := res.Summary()
	if got != "removed 2 key(s), skipped 1 key(s)" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsHelper(s, sub))
}

func containsHelper(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
