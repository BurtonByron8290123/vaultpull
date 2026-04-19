// Package envlock provides advisory file locking for .env files to prevent
// concurrent writes during a pull operation.
package envlock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrLockTimeout is returned when the lock cannot be acquired within the deadline.
var ErrLockTimeout = errors.New("envlock: timed out waiting for lock")

// Locker manages a lock file alongside a target .env file.
type Locker struct {
	lockPath string
	timeout  time.Duration
	pollInterval time.Duration
}

// New returns a Locker for the given env file path.
// The lock file is placed next to the target file with a .lock suffix.
func New(envPath string, timeout time.Duration) *Locker {
	return &Locker{
		lockPath:     envPath + ".lock",
		timeout:      timeout,
		pollInterval: 50 * time.Millisecond,
	}
}

// Acquire creates the lock file, blocking until it is available or the
// timeout elapses. It writes the current PID into the lock file.
func (l *Locker) Acquire() error {
	deadline := time.Now().Add(l.timeout)
	for {
		f, err := os.OpenFile(l.lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err == nil {
			_, _ = fmt.Fprintf(f, "%d", os.Getpid())
			_ = f.Close()
			return nil
		}
		if !os.IsExist(err) {
			return fmt.Errorf("envlock: create lock file: %w", err)
		}
		if time.Now().After(deadline) {
			return ErrLockTimeout
		}
		time.Sleep(l.pollInterval)
	}
}

// Release removes the lock file.
func (l *Locker) Release() error {
	if err := os.Remove(l.lockPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("envlock: remove lock file: %w", err)
	}
	return nil
}

// LockPath returns the resolved path of the lock file.
func (l *Locker) LockPath() string {
	return filepath.Clean(l.lockPath)
}
