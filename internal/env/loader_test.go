package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestParseBasicKeyValue(t *testing.T) {
	p := writeTemp(t, "FOO=bar\nBAZ=qux\n")
	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "FOO" || entries[0].Value != "bar" {
		t.Errorf("entry 0: got %+v", entries[0])
	}
}

func TestParseSkipsCommentsAndBlanks(t *testing.T) {
	p := writeTemp(t, "# comment\n\nKEY=value\n")
	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "KEY" {
		t.Errorf("unexpected key: %s", entries[0].Key)
	}
}

func TestParseUnquotesDoubleQuotes(t *testing.T) {
	p := writeTemp(t, `SECRET="hello world"` + "\n")
	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "hello world" {
		t.Errorf("want 'hello world', got %q", entries[0].Value)
	}
}

func TestParseUnquotesSingleQuotes(t *testing.T) {
	p := writeTemp(t, "TOKEN='abc123'\n")
	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "abc123" {
		t.Errorf("want 'abc123', got %q", entries[0].Value)
	}
}

func TestParseExportPrefix(t *testing.T) {
	p := writeTemp(t, "export MY_VAR=exported\n")
	entries, err := Parse(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "MY_VAR" || entries[0].Value != "exported" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestParseMissingFileReturnsNil(t *testing.T) {
	entries, err := Parse("/nonexistent/.env")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for missing file")
	}
}

func TestToMap(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "A", Value: "3"}, // duplicate — last wins
	}
	m := ToMap(entries)
	if m["A"] != "3" {
		t.Errorf("want A=3, got %s", m["A"])
	}
	if m["B"] != "2" {
		t.Errorf("want B=2, got %s", m["B"])
	}
}
