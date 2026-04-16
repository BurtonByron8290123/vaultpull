// Package metrics provides lightweight run-result tracking for vaultpull.
//
// After each pull operation the caller constructs a RunResult with counts of
// added, updated, removed and unchanged keys together with the elapsed
// duration, then hands it to a Printer to emit a structured summary line.
//
// Example:
//
//	p := metrics.New(os.Stderr)
//	p.Print(metrics.RunResult{
//		Path:      "secret/myapp",
//		Added:     3,
//		Duration:  time.Since(start),
//		Timestamp: time.Now(),
//	})
package metrics
