package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Path      string    `json:"path"`
	Target    string    `json:"target"`
	Added     int       `json:"added"`
	Updated   int       `json:"updated"`
	Removed   int       `json:"removed"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes structured audit entries to a file.
type Logger struct {
	path string
}

// New creates a Logger that appends to the given file path.
func New(path string) *Logger {
	return &Logger{path: path}
}

// Record appends an Entry to the audit log file as a JSON line.
func (l *Logger) Record(e Entry) error {
	if l.path == "" {
		return nil
	}
	e.Timestamp = time.Now().UTC()

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}
