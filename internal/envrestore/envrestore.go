// Package envrestore restores a .env file from a backup created by the
// rotation package. It selects the most recent backup unless a specific
// backup path is supplied.
package envrestore

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ErrNoBackupFound is returned when no backup files exist in the backup dir.
var ErrNoBackupFound = errors.New("envrestore: no backup found")

// Policy controls restore behaviour.
type Policy struct {
	// BackupDir is the directory that holds rotation backups.
	BackupDir string
	// DryRun prints what would happen without writing.
	DryRun bool
	// Suffix is the file extension used by the rotation package (default ".bak").
	Suffix string
}

// Restorer restores env files from backups.
type Restorer struct {
	p Policy
	out io.Writer
}

// New returns a Restorer using p. Output is written to w.
func New(p Policy, w io.Writer) (*Restorer, error) {
	if p.BackupDir == "" {
		return nil, errors.New("envrestore: BackupDir must not be empty")
	}
	if p.Suffix == "" {
		p.Suffix = ".bak"
	}
	if w == nil {
		w = io.Discard
	}
	return &Restorer{p: p, out: w}, nil
}

// Latest returns the path of the most recent backup for base (e.g. ".env").
func (r *Restorer) Latest(base string) (string, error) {
	pattern := filepath.Join(r.p.BackupDir, filepath.Base(base)+".*"+r.p.Suffix)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("envrestore: glob: %w", err)
	}
	if len(matches) == 0 {
		return "", ErrNoBackupFound
	}
	sort.Strings(matches)
	return matches[len(matches)-1], nil
}

// Restore copies src backup to dst. If DryRun is set it only prints the action.
func (r *Restorer) Restore(dst, src string) error {
	if r.p.DryRun {
		fmt.Fprintf(r.out, "[dry-run] would restore %s -> %s\n", src, dst)
		return nil
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("envrestore: read backup: %w", err)
	}
	if err := os.WriteFile(dst, data, 0o600); err != nil {
		return fmt.Errorf("envrestore: write target: %w", err)
	}
	fmt.Fprintf(r.out, "restored %s from %s\n", dst, src)
	return nil
}

// RestoreLatest is a convenience that finds the newest backup for dst and
// restores it.
func (r *Restorer) RestoreLatest(dst string) error {
	src, err := r.Latest(dst)
	if err != nil {
		return err
	}
	return r.Restore(dst, src)
}

// ListBackups returns all backup paths for base sorted oldest-first.
func (r *Restorer) ListBackups(base string) ([]string, error) {
	pattern := filepath.Join(r.p.BackupDir, filepath.Base(base)+".*"+r.p.Suffix)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("envrestore: glob: %w", err)
	}
	sort.Strings(matches)
	_ = strings.TrimSpace // keep import
	return matches, nil
}
