package metrics

import (
	"fmt"
	"io"
	"os"
	"time"
)

// RunResult holds statistics from a single pull run.
type RunResult struct {
	Path      string
	Added     int
	Updated   int
	Removed   int
	Unchanged int
	Duration  time.Duration
	Timestamp time.Time
}

// Summary returns a human-readable one-line summary.
func (r RunResult) Summary() string {
	return fmt.Sprintf(
		"[%s] path=%s added=%d updated=%d removed=%d unchanged=%d duration=%s",
		r.Timestamp.Format(time.RFC3339),
		r.Path,
		r.Added,
		r.Updated,
		r.Removed,
		r.Unchanged,
		r.Duration.Round(time.Millisecond),
	)
}

// Printer writes RunResult summaries to an io.Writer.
type Printer struct {
	w io.Writer
}

// New returns a Printer that writes to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w}
}

// Print writes the summary of r to the underlying writer.
func (p *Printer) Print(r RunResult) {
	fmt.Fprintln(p.w, r.Summary())
}
