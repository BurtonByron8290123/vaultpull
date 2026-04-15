package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry holds a cached secret payload with an expiry timestamp.
type Entry struct {
	Secrets   map[string]string `json:"secrets"`
	FetchedAt time.Time         `json:"fetched_at"`
	TTL       time.Duration     `json:"ttl"`
}

// IsExpired reports whether the cache entry is past its TTL.
func (e *Entry) IsExpired() bool {
	if e.TTL <= 0 {
		return false
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// Store is a file-backed cache for Vault secrets.
type Store struct {
	mu  sync.Mutex
	dir string
}

// NewStore creates a Store that persists cache files under dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) path(key string) string {
	return filepath.Join(s.dir, key+".json")
}

// Get returns the cached Entry for key, or (nil, nil) if not found.
func (s *Store) Get(key string) (*Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path(key))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// Set writes an Entry to disk for key.
func (s *Store) Set(key string, e *Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(s.dir, 0700); err != nil {
		return err
	}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(key), data, 0600)
}

// Invalidate removes the cached entry for key.
func (s *Store) Invalidate(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Remove(s.path(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
