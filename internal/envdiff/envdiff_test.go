package envdiff_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envdiff"
)

func TestCompareAdded(t *testing.T) {
	r := envdiff.Compare(map[string]string{}, map[string]string{"FOO": "bar"})
	if len(r.Changes) != 1 || r.Changes[0].Kind != envdiff.Added {
		t.Fatalf("expected one Added change, got %+v", r.Changes)
	}
}

func TestCompareRemoved(t *testing.T) {
	r := envdiff.Compare(map[string]string{"FOO": "bar"}, map[string]string{})
	if len(r.Changes) != 1 || r.Changes[0].Kind != envdiff.Removed {
		t.Fatalf("expected one Removed change, got %+v", r.Changes)
	}
}

func TestCompareUpdated(t *testing.T) {
	r := envdiff.Compare(map[string]string{"FOO": "old"}, map[string]string{"FOO": "new"})
	if len(r.Changes) != 1 || r.Changes[0].Kind != envdiff.Updated {
		t.Fatalf("expected one Updated change, got %+v", r.Changes)
	}
	if r.Changes[0].Old != "old" || r.Changes[0].New != "new" {
		t.Fatalf("unexpected old/new values: %+v", r.Changes[0])
	}
}

func TestCompareUnchanged(t *testing.T) {
	r := envdiff.Compare(map[string]string{"FOO": "bar"}, map[string]string{"FOO": "bar"})
	if r.HasChanges() {
		t.Fatal("expected no changes")
	}
}

func TestSummaryFormat(t *testing.T) {
	prev := map[string]string{"OLD": "v"}
	next := map[string]string{"NEW": "v", "UPD": "b"}
	prev["UPD"] = "a"
	r := envdiff.Compare(prev, next)
	s := r.Summary()
	if !strings.Contains(s, "+1") || !strings.Contains(s, "~1") || !strings.Contains(s, "-1") {
		t.Fatalf("unexpected summary: %s", s)
	}
}

func TestFprintOutput(t *testing.T) {
	prev := map[string]string{"A": "old"}
	next := map[string]string{"A": "new", "B": "val"}
	r := envdiff.Compare(prev, next)
	var sb strings.Builder
	envdiff.Fprint(&sb, r)
	out := sb.String()
	if !strings.Contains(out, "+ B=") {
		t.Errorf("expected added line, got: %s", out)
	}
	if !strings.Contains(out, "~ A=") {
		t.Errorf("expected updated line, got: %s", out)
	}
}

func TestChangesSortedByKey(t *testing.T) {
	next := map[string]string{"Z": "1", "A": "2", "M": "3"}
	r := envdiff.Compare(map[string]string{}, next)
	keys := make([]string, len(r.Changes))
	for i, c := range r.Changes {
		keys[i] = c.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Fatalf("changes not sorted: %v", keys)
		}
	}
}
