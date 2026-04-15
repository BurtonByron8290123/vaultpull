package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintAdded(t *testing.T) {
	r := &Result{
		Changes: []Change{
			{Key: "API_KEY", Type: Added, NewVal: "[6 chars]"},
		},
	}

	var buf bytes.Buffer
	Print(&buf, r, PrintOptions{NoColor: true})

	if !strings.Contains(buf.String(), "+ API_KEY") {
		t.Errorf("expected '+ API_KEY' in output, got:\n%s", buf.String())
	}
}

func TestPrintRemoved(t *testing.T) {
	r := &Result{
		Changes: []Change{
			{Key: "OLD_KEY", Type: Removed, OldVal: "[3 chars]"},
		},
	}

	var buf bytes.Buffer
	Print(&buf, r, PrintOptions{NoColor: true})

	if !strings.Contains(buf.String(), "- OLD_KEY") {
		t.Errorf("expected '- OLD_KEY' in output, got:\n%s", buf.String())
	}
}

func TestPrintUpdated(t *testing.T) {
	r := &Result{
		Changes: []Change{
			{Key: "DB_PASS", Type: Updated, OldVal: "[3 chars]", NewVal: "[8 chars]"},
		},
	}

	var buf bytes.Buffer
	Print(&buf, r, PrintOptions{NoColor: true})

	if !strings.Contains(buf.String(), "~ DB_PASS") {
		t.Errorf("expected '~ DB_PASS' in output, got:\n%s", buf.String())
	}
}

func TestPrintUnchangedHiddenByDefault(t *testing.T) {
	r := &Result{
		Changes: []Change{
			{Key: "STABLE", Type: Unchanged},
		},
	}

	var buf bytes.Buffer
	Print(&buf, r, PrintOptions{NoColor: true, ShowUnchanged: false})

	if strings.Contains(buf.String(), "STABLE") {
		t.Error("unchanged key should be hidden when ShowUnchanged is false")
	}
}

func TestPrintUnchangedVisibleWhenEnabled(t *testing.T) {
	r := &Result{
		Changes: []Change{
			{Key: "STABLE", Type: Unchanged},
		},
	}

	var buf bytes.Buffer
	Print(&buf, r, PrintOptions{NoColor: true, ShowUnchanged: true})

	if !strings.Contains(buf.String(), "STABLE") {
		t.Error("unchanged key should be visible when ShowUnchanged is true")
	}
}

func TestPrintSummaryAlwaysPresent(t *testing.T) {
	r := &Result{Changes: []Change{}}

	var buf bytes.Buffer
	Print(&buf, r, PrintOptions{NoColor: true})

	if !strings.Contains(buf.String(), "added") {
		t.Error("summary line should always be printed")
	}
}
