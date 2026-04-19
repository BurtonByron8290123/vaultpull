package labelmap_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/labelmap"
)

func TestApplyRenamesKey(t *testing.T) {
	m, err := labelmap.New([]labelmap.Entry{{From: "DB_PASS", To: "DATABASE_PASSWORD"}})
	if err != nil {
		t.Fatal(err)
	}
	out := m.Apply(map[string]string{"DB_PASS": "secret", "OTHER": "val"})
	if _, ok := out["DB_PASS"]; ok {
		t.Error("old key should be removed")
	}
	if out["DATABASE_PASSWORD"] != "secret" {
		t.Errorf("expected DATABASE_PASSWORD=secret, got %q", out["DATABASE_PASSWORD"])
	}
	if out["OTHER"] != "val" {
		t.Error("unmatched key should be preserved")
	}
}

func TestApplyNoMatchPassesThrough(t *testing.T) {
	m, _ := labelmap.New([]labelmap.Entry{{From: "X", To: "Y"}})
	out := m.Apply(map[string]string{"A": "1"})
	if out["A"] != "1" {
		t.Error("key should be unchanged")
	}
}

func TestNewRejectsEmptyFrom(t *testing.T) {
	_, err := labelmap.New([]labelmap.Entry{{From: "", To: "B"}})
	if err == nil {
		t.Error("expected error for empty from")
	}
}

func TestNewRejectsEmptyTo(t *testing.T) {
	_, err := labelmap.New([]labelmap.Entry{{From: "A", To: ""}})
	if err == nil {
		t.Error("expected error for empty to")
	}
}

func TestLoadConfigValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.yaml")
	content := "- from: SECRET_KEY\n  to: APP_SECRET\n"
	os.WriteFile(path, []byte(content), 0o600)

	m, err := labelmap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	out := m.Apply(map[string]string{"SECRET_KEY": "abc"})
	if out["APP_SECRET"] != "abc" {
		t.Errorf("unexpected value: %v", out)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := labelmap.LoadConfig("/nonexistent/labels.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestApplyMultipleRules(t *testing.T) {
	m, _ := labelmap.New([]labelmap.Entry{
		{From: "A", To: "AA"},
		{From: "B", To: "BB"},
	})
	out := m.Apply(map[string]string{"A": "1", "B": "2", "C": "3"})
	if out["AA"] != "1" || out["BB"] != "2" || out["C"] != "3" {
		t.Errorf("unexpected result: %v", out)
	}
}
