package planmode_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/diff"
	"github.com/your-org/vaultpull/internal/planmode"
)

func TestBuildDetectsAdded(t *testing.T) {
	p := planmode.Build(".env", map[string]string{}, map[string]string{"FOO": "bar"})
	if len(p.Changes) != 1 || p.Changes[0].Op != diff.Added {
		t.Fatalf("expected one Added change, got %+v", p.Changes)
	}
}

func TestBuildDetectsUpdated(t *testing.T) {
	p := planmode.Build(".env", map[string]string{"FOO": "old"}, map[string]string{"FOO": "new"})
	if len(p.Changes) != 1 || p.Changes[0].Op != diff.Updated {
		t.Fatalf("expected one Updated change, got %+v", p.Changes)
	}
}

func TestBuildDetectsRemoved(t *testing.T) {
	p := planmode.Build(".env", map[string]string{"FOO": "bar"}, map[string]string{})
	if len(p.Changes) != 1 || p.Changes[0].Op != diff.Removed {
		t.Fatalf("expected one Removed change, got %+v", p.Changes)
	}
}

func TestBuildUnchanged(t *testing.T) {
	p := planmode.Build(".env", map[string]string{"FOO": "bar"}, map[string]string{"FOO": "bar"})
	if p.HasChanges() {
		t.Fatal("expected no changes")
	}
}

func TestSummaryFormat(t *testing.T) {
	p := planmode.Build(".env",
		map[string]string{"OLD": "x", "SAME": "y"},
		map[string]string{"NEW": "a", "SAME": "z"},
	)
	s := p.Summary()
	if !strings.Contains(s, "+1") || !strings.Contains(s, "~1") || !strings.Contains(s, "-1") {
		t.Fatalf("unexpected summary: %s", s)
	}
}

func TestPrintShowsNoChanges(t *testing.T) {
	p := planmode.Build(".env", map[string]string{"K": "v"}, map[string]string{"K": "v"})
	var buf bytes.Buffer
	planmode.Print(p, &buf)
	if !strings.Contains(buf.String(), "no changes") {
		t.Fatalf("expected 'no changes' in output, got: %s", buf.String())
	}
}

func TestPrintShowsOperationSymbols(t *testing.T) {
	p := planmode.Build(".env",
		map[string]string{"B": "old"},
		map[string]string{"A": "new", "B": "updated"},
	)
	var buf bytes.Buffer
	planmode.Print(p, &buf)
	out := buf.String()
	if !strings.Contains(out, "+ A") {
		t.Errorf("expected '+ A' in output")
	}
	if !strings.Contains(out, "~ B") {
		t.Errorf("expected '~ B' in output")
	}
}
