// Package envaudit provides a policy-driven audit trail for env file changes,
// recording which keys were added, updated, or removed along with timestamps.
package envaudit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ChangeKind describes the type of change recorded.
type ChangeKind string

const (
	ChangeAdded   ChangeKind = "added"
	ChangeUpdated ChangeKind = "updated"
	ChangeRemoved ChangeKind = "removed"
)

// Entry is a single audit record.
type Entry struct {
	Timestamp time.Time  `json:"timestamp"`
	Path      string     `json:"path"`
	Key       string     `json:"key"`
	Kind      ChangeKind `json:"kind"`
	Masked    bool       `json:"masked"`
}

// Policy controls audit behaviour.
type Policy struct {
	// AuditPath is the file where entries are appended. Empty disables auditing.
	AuditPath string
	// MaskValues redacts secret values from log entries.
	MaskValues bool
}

// Logger writes audit entries to a JSONL file.
type Logger struct {
	policy Policy
	clock  func() time.Time
}

// New returns a Logger using the given policy.
func New(p Policy) *Logger {
	return &Logger{policy: p, clock: time.Now}
}

// Record appends one entry per changed key to the audit log.
func (l *Logger) Record(envPath string, changes map[string]ChangeKind) error {
	if l.policy.AuditPath == "" || len(changes) == 0 {
		return nil
	}
	f, err := os.OpenFile(l.policy.AuditPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("envaudit: open %s: %w", l.policy.AuditPath, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for key, kind := range changes {
		e := Entry{
			Timestamp: l.clock().UTC(),
			Path:      envPath,
			Key:       key,
			Kind:      kind,
			Masked:    l.policy.MaskValues,
		}
		if err := enc.Encode(e); err != nil {
			return fmt.Errorf("envaudit: encode entry: %w", err)
		}
	}
	return nil
}
