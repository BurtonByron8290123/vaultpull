// Package multienv supports writing secrets to multiple .env files
// based on environment-specific mappings (e.g. .env.dev, .env.prod).
package multienv

import (
	"fmt"
	"path/filepath"

	"github.com/your-org/vaultpull/internal/env"
)

// Target describes a single output .env file and which keys to include.
type Target struct {
	// Name is a human-readable label (e.g. "dev", "prod").
	Name string `yaml:"name"`
	// Path is the file path to write (e.g. ".env.dev").
	Path string `yaml:"path"`
	// Keys lists the secret keys to include. Empty means all keys.
	Keys []string `yaml:"keys,omitempty"`
}

// Writer writes secrets to multiple target .env files.
type Writer struct {
	targets []Target
	dir     string
}

// New creates a Writer for the given targets rooted at dir.
func New(dir string, targets []Target) *Writer {
	return &Writer{targets: targets, dir: dir}
}

// WriteAll writes the provided secrets to each configured target file.
// Keys not listed in a target's Keys slice are omitted; an empty Keys
// slice means all keys are written.
func (w *Writer) WriteAll(secrets map[string]string) error {
	for _, t := range w.targets {
		filtered := filter(secrets, t.Keys)
		path := filepath.Join(w.dir, t.Path)
		if err := env.WriteFile(path, filtered, 0o600); err != nil {
			return fmt.Errorf("multienv: write target %q (%s): %w", t.Name, path, err)
		}
	}
	return nil
}

// filter returns a subset of src containing only the specified keys.
// If keys is empty, src is returned as-is.
func filter(src map[string]string, keys []string) map[string]string {
	if len(keys) == 0 {
		return src
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := src[k]; ok {
			out[k] = v
		}
	}
	return out
}
