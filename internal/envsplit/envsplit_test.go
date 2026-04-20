package envsplit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envsplit"
)

func TestApplySplitsByPrefix(t *testing.T) {
	p := envsplit.Policy{
		Rules: []envsplit.Rule{
			{Prefix: "APP_", Group: "app"},
			{Prefix: "DB_", Group: "database"},
		},
	}
	s, err := envsplit.New(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	src := map[string]string{
		"APP_HOST": "localhost",
		"DB_URL":   "postgres://",
		"OTHER":    "value",
	}
	out := s.Apply(src)
	if out["app"]["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST in app group")
	}
	if out["database"]["DB_URL"] != "postgres://" {
		t.Errorf("expected DB_URL in database group")
	}
	if _, ok := out["app"]["OTHER"]; ok {
		t.Errorf("OTHER should not appear in app group")
	}
}

func TestApplyStripsPrefix(t *testing.T) {
	p := envsplit.Policy{
		Rules: []envsplit.Rule{
			{Prefix: "APP_", Group: "app", Strip: true},
		},
	}
	s, _ := envsplit.New(p)
	out := s.Apply(map[string]string{"APP_PORT": "8080"})
	if out["app"]["PORT"] != "8080" {
		t.Errorf("expected stripped key PORT, got %v", out["app"])
	}
}

func TestApplyDefaultGroup(t *testing.T) {
	p := envsplit.Policy{
		Rules:        []envsplit.Rule{{Prefix: "APP_", Group: "app"}},
		DefaultGroup: "misc",
	}
	s, _ := envsplit.New(p)
	out := s.Apply(map[string]string{"UNMATCHED": "yes"})
	if out["misc"]["UNMATCHED"] != "yes" {
		t.Errorf("expected UNMATCHED in misc group")
	}
}

func TestNewRejectsEmptyPrefix(t *testing.T) {
	p := envsplit.Policy{
		Rules: []envsplit.Rule{{Prefix: "", Group: "app"}},
	}
	_, err := envsplit.New(p)
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNewRejectsEmptyGroup(t *testing.T) {
	p := envsplit.Policy{
		Rules: []envsplit.Rule{{Prefix: "APP_", Group: ""}},
	}
	_, err := envsplit.New(p)
	if err == nil {
		t.Fatal("expected error for empty group")
	}
}

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "split.yaml")
	content := []byte(`rules:\n  - prefix: SVC_\n    group: service\n    strip: true\ndefault_group: other\n`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := envsplit.LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := envsplit.LoadConfig("/nonexistent/split.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
