package envcompare_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/vaultpull/internal/envcompare"
)

func writeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCompareAdded(t *testing.T) {
	src := map[string]string{}
	dst := map[string]string{"FOO": "bar"}
	r := envcompare.Compare(src, dst)
	if len(r.Changes) != 1 || r.Changes[0].Kind != envcompare.Added {
		t.Fatalf("expected 1 added change, got %+v", r.Changes)
	}
}

func TestCompareRemoved(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	dst := map[string]string{}
	r := envcompare.Compare(src, dst)
	if len(r.Changes) != 1 || r.Changes[0].Kind != envcompare.Removed {
		t.Fatalf("expected 1 removed change, got %+v", r.Changes)
	}
}

func TestCompareUpdated(t *testing.T) {
	src := map[string]string{"FOO": "old"}
	dst := map[string]string{"FOO": "new"}
	r := envcompare.Compare(src, dst)
	if len(r.Changes) != 1 || r.Changes[0].Kind != envcompare.Updated {
		t.Fatalf("expected 1 updated change, got %+v", r.Changes)
	}
	if r.Changes[0].OldValue != "old" || r.Changes[0].NewValue != "new" {
		t.Fatalf("unexpected values: %+v", r.Changes[0])
	}
}

func TestCompareUnchanged(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	dst := map[string]string{"FOO": "bar"}
	r := envcompare.Compare(src, dst)
	if len(r.Changes) != 0 {
		t.Fatalf("expected no changes, got %+v", r.Changes)
	}
}

func TestSummaryFormat(t *testing.T) {
	r := &envcompare.Result{
		Changes: []envcompare.Change{
			{Key: "A", Kind: envcompare.Added},
			{Key: "B", Kind: envcompare.Updated},
			{Key: "C", Kind: envcompare.Removed},
		},
	}
	got := r.Summary()
	expected := "+1 added, ~1 updated, -1 removed"
	if got != expected {
		t.Fatalf("expected %q got %q", expected, got)
	}
}

func TestFprintOutput(t *testing.T) {
	r := envcompare.Compare(
		map[string]string{"OLD": "v1"},
		map[string]string{"NEW": "v2"},
	)
	var buf bytes.Buffer
	envcompare.Fprint(&buf, r)
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestConfigMaskApply(t *testing.T) {
	cfg := envcompare.DefaultConfig()
	r := &envcompare.Result{
		Changes: []envcompare.Change{
			{Key: "SECRET", Kind: envcompare.Updated, OldValue: "real", NewValue: "also-real"},
		},
	}
	masked := cfg.Apply(r)
	if masked.Changes[0].NewValue != "***" {
		t.Fatalf("expected masked value, got %q", masked.Changes[0].NewValue)
	}
}

func TestCompareFiles(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, "src.env", "FOO=bar\nBAZ=qux\n")
	dst := writeEnv(t, dir, "dst.env", "FOO=bar\nNEW=val\n")
	r, err := envcompare.CompareFiles(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Changes) != 2 {
		t.Fatalf("expected 2 changes (added+removed), got %d: %+v", len(r.Changes), r.Changes)
	}
}
