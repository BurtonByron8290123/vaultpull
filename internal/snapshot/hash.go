package snapshot

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"time"
)

// HashValue returns a short SHA-256 hex digest of the given secret value.
// Only the first 8 characters are kept to avoid storing sensitive material.
func HashValue(value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", sum[:4])
}

// Build constructs a Snapshot from a map of key→value pairs.
func Build(vaultPath string, secrets map[string]string) Snapshot {
	entries := make([]Entry, 0, len(secrets))
	now := time.Now().UTC()

	for k, v := range secrets {
		entries = append(entries, Entry{
			Key:        k,
			ValueHash:  HashValue(v),
			CapturedAt: now,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Snapshot{
		Path:       vaultPath,
		CapturedAt: now,
		Entries:    entries,
	}
}
