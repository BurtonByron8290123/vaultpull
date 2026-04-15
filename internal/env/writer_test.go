package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMarshalBasic(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	data, err := Marshal(vars, MarshalOptions{SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := string(data)
	if !strings.Contains(got, "BAZ=qux\n") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", got)
	}
	if !strings.Contains(got, "FOO=bar\n") {
		t.Errorf("expected FOO=bar in output, got:\n%s", got)
	}
}

func TestMarshalSortedOrder(t *testing.T) {
	vars := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	data, err := Marshal(vars, MarshalOptions{SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line to be A_KEY, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected last line to be Z_KEY, got %q", lines[2])
	}
}

func TestMarshalQuotesValuesWithSpaces(t *testing.T) {
	vars := map[string]string{"MSG": "hello world"}
	data, err := Marshal(vars, MarshalOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(data), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(data))
	}
}

func TestMarshalInvalidKeyReturnsError(t *testing.T) {
	vars := map[string]string{"INVALID-KEY": "value"}
	_, err := Marshal(vars, MarshalOptions{})
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

func TestMarshalEmptyKeyReturnsError(t *testing.T) {
	vars := map[string]string{"": "value"}
	_, err := Marshal(vars, MarshalOptions{})
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestWriteFileCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	vars := map[string]string{"TOKEN": "secret123"}
	err := WriteFile(path, vars, MarshalOptions{}, 0600)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), "TOKEN=secret123") {
		t.Errorf("expected TOKEN=secret123 in file, got: %s", string(data))
	}
}

func TestWriteFilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	err := WriteFile(path, map[string]string{"K": "v"}, MarshalOptions{}, 0600)
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected perm 0600, got %v", info.Mode().Perm())
	}
}
