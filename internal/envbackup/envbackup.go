// Package envbackup provides utilities for creating and managing timestamped
// backups of .env files before mutations are applied.
package envbackup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		BackupDir:  ".env-backups",
		Suffix:     ".bak",
		MaxBackups: 10,
		DryRun:     false,
	}
}

// Policy controls backup behaviour.
type Policy struct {
	// BackupDir is the directory where backup files are written.
	BackupDir string

	// Suffix is appended to each backup filename before the timestamp.
	Suffix string

	// MaxBackups is the maximum number of backups to retain per source file.
	// Zero means unlimited.
	MaxBackups int

	// DryRun skips all writes when true.
	DryRun bool
}

// Manager creates and prunes backups for a single .env file.
type Manager struct {
	policy Policy
	clock  func() time.Time
}

// New returns a Manager using the given policy.
func New(p Policy) (*Manager, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	return &Manager{policy: p, clock: time.Now}, nil
}

func validate(p Policy) error {
	if strings.TrimSpace(p.BackupDir) == "" {
		return fmt.Errorf("envbackup: BackupDir must not be empty")
	}
	if p.MaxBackups < 0 {
		return fmt.Errorf("envbackup: MaxBackups must be >= 0")
	}
	return nil
}

// Backup copies src into BackupDir with a timestamp-based name and then prunes
// old backups according to MaxBackups. It returns the path of the new backup
// file, or an empty string when DryRun is true.
func (m *Manager) Backup(src string) (string, error) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		// Nothing to back up.
		return "", nil
	}

	if m.policy.DryRun {
		return "", nil
	}

	if err := os.MkdirAll(m.policy.BackupDir, 0o700); err != nil {
		return "", fmt.Errorf("envbackup: create backup dir: %w", err)
	}

	base := filepath.Base(src)
	ts := m.clock().UTC().Format("20060102T150405Z")
	dst := filepath.Join(m.policy.BackupDir, base+m.policy.Suffix+"."+ts)

	if err := copyFile(src, dst); err != nil {
		return "", fmt.Errorf("envbackup: copy %s -> %s: %w", src, dst, err)
	}

	if err := m.prune(base); err != nil {
		return dst, fmt.Errorf("envbackup: prune: %w", err)
	}

	return dst, nil
}

// prune removes old backups for the given base filename, keeping at most
// MaxBackups entries. When MaxBackups is 0 nothing is pruned.
func (m *Manager) prune(base string) error {
	if m.policy.MaxBackups == 0 {
		return nil
	}

	pattern := filepath.Join(m.policy.BackupDir, base+m.policy.Suffix+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	sort.Strings(matches) // ascending timestamp order

	for len(matches) > m.policy.MaxBackups {
		if err := os.Remove(matches[0]); err != nil && !os.IsNotExist(err) {
			return err
		}
		matches = matches[1:]
	}
	return nil
}

// List returns all backup paths for the given source filename in ascending
// chronological order.
func (m *Manager) List(src string) ([]string, error) {
	base := filepath.Base(src)
	pattern := filepath.Join(m.policy.BackupDir, base+m.policy.Suffix+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)
	return matches, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
