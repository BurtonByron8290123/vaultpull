package envwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := New(path, true)
	err := w.Write(map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output, got:\n%s", content)
	}
}

func TestWriteOverwriteReplacesKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	_ = os.WriteFile(path, []byte("DB_HOST=oldhost\nAPI_KEY=secret\n"), 0600)

	w := New(path, true)
	err := w.Write(map[string]string{"DB_HOST": "newhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if strings.Contains(content, "API_KEY") {
		t.Errorf("expected API_KEY to be removed in overwrite mode")
	}
	if !strings.Contains(content, "DB_HOST=newhost") {
		t.Errorf("expected DB_HOST=newhost, got:\n%s", content)
	}
}

func TestWriteMergePreservesExistingKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	_ = os.WriteFile(path, []byte("EXISTING_KEY=preserved\n"), 0600)

	w := New(path, false)
	err := w.Write(map[string]string{"NEW_KEY": "added"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "EXISTING_KEY=preserved") {
		t.Errorf("expected EXISTING_KEY to be preserved, got:\n%s", content)
	}
	if !strings.Contains(content, "NEW_KEY=added") {
		t.Errorf("expected NEW_KEY=added to be present, got:\n%s", content)
	}
}

func TestWriteFilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := New(path, true)
	_ = w.Write(map[string]string{"SECRET": "value"})

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file permissions 0600, got %v", info.Mode().Perm())
	}
}
