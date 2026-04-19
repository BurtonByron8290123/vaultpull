// Package envpin records the Vault secret versions that were last pulled
// and warns when a secret has changed since it was pinned.
package envpin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// ErrPinNotFound is returned when no pin file exists for the given path.
var ErrPinNotFound = errors.New("envpin: pin file not found")

// Entry records the version metadata for a single secret path.
type Entry struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	PinnedAt  time.Time `json:"pinned_at"`
	Checksum  string    `json:"checksum"`
}

// Store persists pin entries to a JSON file.
type Store struct {
	filePath string
}

// NewStore creates a Store backed by filePath.
func NewStore(filePath string) *Store {
	return &Store{filePath: filePath}
}

// Save writes all entries to the pin file, replacing any existing content.
func (s *Store) Save(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("envpin: marshal: %w", err)
	}
	if err := os.WriteFile(s.filePath, data, 0o600); err != nil {
		return fmt.Errorf("envpin: write: %w", err)
	}
	return nil
}

// Load reads all entries from the pin file.
// Returns ErrPinNotFound if the file does not exist.
func (s *Store) Load() ([]Entry, error) {
	data, err := os.ReadFile(s.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrPinNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("envpin: read: %w", err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("envpin: unmarshal: %w", err)
	}
	return entries, nil
}

// Diff compares current entries against pinned entries and returns paths
// whose checksum has changed since they were pinned.
func Diff(pinned, current []Entry) []string {
	index := make(map[string]string, len(pinned))
	for _, e := range pinned {
		index[e.Path] = e.Checksum
	}
	var drifted []string
	for _, c := range current {
		if prev, ok := index[c.Path]; ok && prev != c.Checksum {
			drifted = append(drifted, c.Path)
		}
	}
	return drifted
}
