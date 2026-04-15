package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// Entry represents a single secret snapshot captured from Vault.
type Entry struct {
	Key       string    `json:"key"`
	ValueHash string    `json:"value_hash"`
	CapturedAt time.Time `json:"captured_at"`
}

// Snapshot holds all secret entries captured at a point in time.
type Snapshot struct {
	Path      string    `json:"vault_path"`
	CapturedAt time.Time `json:"captured_at"`
	Entries   []Entry   `json:"entries"`
}

// Store persists and retrieves snapshots from disk.
type Store struct {
	filePath string
}

// NewStore creates a new Store backed by the given file path.
func NewStore(filePath string) *Store {
	return &Store{filePath: filePath}
}

// Save writes the snapshot to disk as JSON.
func (s *Store) Save(snap Snapshot) error {
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0600)
}

// Load reads the most recent snapshot from disk.
// Returns nil, nil if no snapshot exists yet.
func (s *Store) Load() (*Snapshot, error) {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}
