package filter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/filter"
)

func TestResolveInlinePatterns(t *testing.T) {
	pc := filter.PatternConfig{
		Patterns: []string{"APP_"},
	}
	f, err := pc.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := f.Apply(map[string]string{
		"APP_KEY": "val",
		"DB_HOST": "host",
	})
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
}

func TestResolveFromFile(t *testing.T) {
	dir := t.TempDir()
	pf := filepath.Join(dir, "patterns.txt")
	content := "# comment\nDB_\n\nAPP_\n"
	if err := os.WriteFile(pf, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	pc := filter.PatternConfig{File: pf}
	f, err := pc.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := f.Apply(map[string]string{
		"APP_KEY": "val",
		"DB_HOST": "host",
		"OTHER":   "x",
	})
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestResolveMissingFileReturnsError(t *testing.T) {
	pc := filter.PatternConfig{File: "/nonexistent/patterns.txt"}
	_, err := pc.Resolve()
	if err == nil {
		t.Error("expected error for missing pattern file")
	}
}

func TestResolveInlineAndFileMerged(t *testing.T) {
	dir := t.TempDir()
	pf := filepath.Join(dir, "patterns.txt")
	if err := os.WriteFile(pf, []byte("!DB_PASSWORD\n"), 0600); err != nil {
		t.Fatal(err)
	}

	pc := filter.PatternConfig{
		Patterns: []string{"DB_"},
		File:     pf,
	}
	f, err := pc.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := f.Apply(map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
	})
	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("expected DB_PASSWORD to be excluded")
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST}
}
