package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/metrics"
)

func result() metrics.RunResult {
	return metrics.RunResult{
		Path:      "secret/app",
		Added:     2,
		Updated:   1,
		Removed:   0,
		Unchanged: 5,
		Duration:  150 * time.Millisecond,
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}
}

func TestSummaryContainsPath(t *testing.T) {
	r := result()
	if !strings.Contains(r.Summary(), "secret/app") {
		t.Fatalf("expected path in summary, got: %s", r.Summary())
	}
}

func TestSummaryContainsCounts(t *testing.T) {
	s := result().Summary()
	for _, want := range []string{"added=2", "updated=1", "removed=0", "unchanged=5"} {
		if !strings.Contains(s, want) {
			t.Errorf("expected %q in summary %q", want, s)
		}
	}
}

func TestSummaryContainsDuration(t *testing.T) {
	s := result().Summary()
	if !strings.Contains(s, "150ms") {
		t.Fatalf("expected duration in summary, got: %s", s)
	}
}

func TestPrintWritesToWriter(t *testing.T) {
	var buf bytes.Buffer
	p := metrics.New(&buf)
	p.Print(result())
	if buf.Len() == 0 {
		t.Fatal("expected output written to buffer")
	}
	if !strings.Contains(buf.String(), "secret/app") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestPrintNewlineTerminated(t *testing.T) {
	var buf bytes.Buffer
	metrics.New(&buf).Print(result())
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Fatal("expected output to end with newline")
	}
}
