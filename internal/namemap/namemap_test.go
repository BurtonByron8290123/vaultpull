package namemap_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/namemap"
)

func rules(t *testing.T, pairs ...string) []namemap.Rule {
	t.Helper()
	var out []namemap.Rule
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, namemap.Rule{From: pairs[i], To: pairs[i+1]})
	}
	return out
}

func TestApplyRenamesKey(t *testing.T) {
	m, err := namemap.New(rules(t, "DB_PASS", "DATABASE_PASSWORD"))
	if err != nil {
		t.Fatal(err)
	}
	if got := m.Apply("DB_PASS"); got != "DATABASE_PASSWORD" {
		t.Fatalf("want DATABASE_PASSWORD, got %s", got)
	}
}

func TestApplyNoMatchReturnsOriginal(t *testing.T) {
	m, _ := namemap.New(rules(t, "A", "B"))
	if got := m.Apply("X"); got != "X" {
		t.Fatalf("want X, got %s", got)
	}
}

func TestApplyMapRewrites(t *testing.T) {
	m, _ := namemap.New(rules(t, "OLD_KEY", "NEW_KEY"))
	src := map[string]string{"OLD_KEY": "val", "KEEP": "same"}
	out := m.ApplyMap(src)
	if out["NEW_KEY"] != "val" {
		t.Fatalf("expected NEW_KEY=val, got %v", out)
	}
	if out["KEEP"] != "same" {
		t.Fatalf("expected KEEP=same, got %v", out)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Fatal("OLD_KEY should have been removed")
	}
}

func TestNewRejectsEmptyFrom(t *testing.T) {
	_, err := namemap.New([]namemap.Rule{{From: "", To: "B"}})
	if err == nil {
		t.Fatal("expected error for empty from")
	}
}

func TestNewRejectsEmptyTo(t *testing.T) {
	_, err := namemap.New([]namemap.Rule{{From: "A", To: ""}})
	if err == nil {
		t.Fatal("expected error for empty to")
	}
}

func TestLoadConfigValid(t *testing.T) {
	data, _ := json.Marshal([]namemap.Rule{{From: "FOO", To: "BAR"}})
	dir := t.TempDir()
	p := filepath.Join(dir, "map.json")
	os.WriteFile(p, data, 0o600)

	m, err := namemap.LoadConfig(p)
	if err != nil {
		t.Fatal(err)
	}
	if got := m.Apply("FOO"); got != "BAR" {
		t.Fatalf("want BAR, got %s", got)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := namemap.LoadConfig("/nonexistent/map.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
