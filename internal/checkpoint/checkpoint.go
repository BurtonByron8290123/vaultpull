// Package checkpoint tracks the last successful pull timestamp per path,
// allowing vaultpull to skip unnecessary syncs when secrets are unchanged.
package checkpoint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry records the last successful pull for a Vault path.
type Entry struct {
	Path      string    `json:"path"`
	PulledAt  time.Time `json:"pulled_at"`
	SecretHash string   `json:"secret_hash"`
}

// Store persists checkpoint entries to a JSON file.
type Store struct {
	mu      sync.Mutex
	file    string
	entries map[string]Entry
}

// NewStore opens or creates a checkpoint store at the given file path.
func NewStore(file string) (*Store, error) {
	s := &Store{file: file, entries: make(map[string]Entry)}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns the checkpoint entry for a Vault path, and whether it exists.
func (s *Store) Get(path string) (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[path]
	return e, ok
}

// Set records a successful pull for the given path and hash, then persists.
func (s *Store) Set(path, hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[path] = Entry{Path: path, PulledAt: time.Now().UTC(), SecretHash: hash}
	return s.save()
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.file)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.entries)
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.file), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.file, data, 0o600)
}
