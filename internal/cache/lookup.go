package cache

import (
	"fmt"
	"strings"
	"time"
)

// Fetcher is a function that retrieves secrets from Vault for a given path.
type Fetcher func(path string) (map[string]string, error)

// Lookup returns secrets for path, using the cache when the entry is present
// and not expired. If the entry is missing or expired, fn is called and the
// result is stored with the given TTL.
func Lookup(store *Store, path string, ttl time.Duration, fn Fetcher) (map[string]string, error) {
	key := cacheKey(path)

	entry, err := store.Get(key)
	if err != nil {
		return nil, fmt.Errorf("cache get: %w", err)
	}

	if entry != nil && !entry.IsExpired() {
		return entry.Secrets, nil
	}

	secrets, err := fn(path)
	if err != nil {
		return nil, err
	}

	newEntry := &Entry{
		Secrets:   secrets,
		FetchedAt: time.Now(),
		TTL:       ttl,
	}
	if setErr := store.Set(key, newEntry); setErr != nil {
		// Non-fatal: return secrets even if caching fails.
		fmt.Printf("warning: cache write failed: %v\n", setErr)
	}
	return secrets, nil
}

// cacheKey converts a Vault path into a safe filesystem key.
func cacheKey(path string) string {
	r := strings.NewReplacer("/", "_", "\\", "_")
	return r.Replace(strings.Trim(path, "/"))
}
