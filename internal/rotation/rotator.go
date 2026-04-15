package rotation

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Rotator handles backup and rotation of existing .env files before overwriting.
type Rotator struct {
	BackupDir string
	MaxBackups int
}

// New creates a new Rotator with the given backup directory and max backup count.
func New(backupDir string, maxBackups int) *Rotator {
	if backupDir == "" {
		backupDir = ".env_backups"
	}
	if maxBackups <= 0 {
		maxBackups = 5
	}
	return &Rotator{
		BackupDir:  backupDir,
		MaxBackups: maxBackups,
	}
}

// Rotate creates a timestamped backup of the given file if it exists.
// It also prunes old backups beyond MaxBackups.
func (r *Rotator) Rotate(envFilePath string) error {
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.MkdirAll(r.BackupDir, 0700); err != nil {
		return fmt.Errorf("rotation: failed to create backup dir: %w", err)
	}

	baseName := filepath.Base(envFilePath)
	timestamp := time.Now().Format("20060102T150405")
	backupName := fmt.Sprintf("%s.%s.bak", baseName, timestamp)
	backupPath := filepath.Join(r.BackupDir, backupName)

	data, err := os.ReadFile(envFilePath)
	if err != nil {
		return fmt.Errorf("rotation: failed to read source file: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return fmt.Errorf("rotation: failed to write backup file: %w", err)
	}

	return r.pruneBackups(baseName)
}

// pruneBackups removes the oldest backups if count exceeds MaxBackups.
func (r *Rotator) pruneBackups(baseName string) error {
	pattern := filepath.Join(r.BackupDir, baseName+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("rotation: failed to list backups: %w", err)
	}

	if len(matches) <= r.MaxBackups {
		return nil
	}

	// matches are sorted lexicographically; oldest timestamps come first
	toRemove := matches[:len(matches)-r.MaxBackups]
	for _, f := range toRemove {
		if err := os.Remove(f); err != nil {
			return fmt.Errorf("rotation: failed to remove old backup %s: %w", f, err)
		}
	}
	return nil
}
