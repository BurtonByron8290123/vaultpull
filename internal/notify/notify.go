package notify

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Notifier sends human-readable messages to an output writer.
type Notifier struct {
	out    io.Writer
	prefix string
	quiet  bool
}

// New returns a Notifier writing to os.Stderr.
func New(prefix string, quiet bool) *Notifier {
	return &Notifier{out: os.Stderr, prefix: prefix, quiet: quiet}
}

// NewWithWriter returns a Notifier writing to the provided writer.
func NewWithWriter(w io.Writer, prefix string, quiet bool) *Notifier {
	return &Notifier{out: w, prefix: prefix, quiet: quiet}
}

// Info emits an informational message unless quiet mode is enabled.
func (n *Notifier) Info(msg string) {
	if n.quiet {
		return
	}
	n.emit(LevelInfo, msg)
}

// Warn always emits a warning message.
func (n *Notifier) Warn(msg string) {
	n.emit(LevelWarn, msg)
}

// Error always emits an error message.
func (n *Notifier) Error(msg string) {
	n.emit(LevelError, msg)
}

// Infof formats and emits an informational message.
func (n *Notifier) Infof(format string, args ...any) {
	n.Info(fmt.Sprintf(format, args...))
}

// Warnf formats and emits a warning message.
func (n *Notifier) Warnf(format string, args ...any) {
	n.Warn(fmt.Sprintf(format, args...))
}

// Errorf formats and emits an error message.
func (n *Notifier) Errorf(format string, args ...any) {
	n.Error(fmt.Sprintf(format, args...))
}

func (n *Notifier) emit(level Level, msg string) {
	parts := []string{}
	if n.prefix != "" {
		parts = append(parts, fmt.Sprintf("[%s]", n.prefix))
	}
	parts = append(parts, fmt.Sprintf("[%s]", strings.ToUpper(string(level))))
	parts = append(parts, msg)
	fmt.Fprintln(n.out, strings.Join(parts, " "))
}
