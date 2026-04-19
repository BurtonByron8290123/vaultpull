// Package rollback provides functionality to restore .env files from backups
// created by the rotation subsystem.
package rollback

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Store manages rollback operations against a backup directory.
type Store struct {
	backupDir string
}

// New returns a Store rooted at backupDir.
func New(backupDir string) (*Store, error) {
	if strings.TrimSpace(backupDir) == "" {
		return nil, errors.New("rollback: backupDir must not be empty")
	}
	return &Store{backupDir: backupDir}, nil
}

// Latest returns the path of the most recent backup for the given base filename,
// or an empty string when no backups exist.
func (s *Store) Latest(base string) (string, error) {
	entries, err := s.list(base)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return "", nil
	}
	return entries[len(entries)-1], nil
}

// Restore copies the most recent backup for base to dst, overwriting dst.
func (s *Store) Restore(base, dst string) error {
	src, err := s.Latest(base)
	if err != nil {
		return err
	}
	if src == "" {
		return fmt.Errorf("rollback: no backup found for %q", base)
	}
	return copyFile(src, dst)
}

// List returns all backup paths for base sorted oldest-first.
func (s *Store) List(base string) ([]string, error) {
	return s.list(base)
}

func (s *Store) list(base string) ([]string, error) {
	pattern := filepath.Join(s.backupDir, base+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("rollback: glob: %w", err)
	}
	sort.Strings(matches)
	return matches, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("rollback: open src: %w", err)
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("rollback: open dst: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("rollback: copy: %w", err)
	}
	return out.Close()
}
