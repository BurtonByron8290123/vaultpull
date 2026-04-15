package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/cache"
)

func TestLookupCallsFetcherOnCacheMiss(t *testing.T) {
	store := cache.NewStore(tempDir(t))
	calls := 0
	fetcher := func(path string) (map[string]string, error) {
		calls++
		return map[string]string{"KEY": "value"}, nil
	}

	secrets, err := cache.Lookup(store, "secret/app", time.Minute, fetcher)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", secrets["KEY"])
	}
	if calls != 1 {
		t.Errorf("expected 1 fetch call, got %d", calls)
	}
}

func TestLookupUsesCacheOnHit(t *testing.T) {
	store := cache.NewStore(tempDir(t))
	calls := 0
	fetcher := func(path string) (map[string]string, error) {
		calls++
		return map[string]string{"KEY": "cached"}, nil
	}

	// Populate cache.
	_, _ = cache.Lookup(store, "secret/app", time.Minute, fetcher)
	// Second call should hit cache.
	_, _ = cache.Lookup(store, "secret/app", time.Minute, fetcher)

	if calls != 1 {
		t.Errorf("expected 1 fetch call, got %d", calls)
	}
}

func TestLookupRefetchesWhenExpired(t *testing.T) {
	store := cache.NewStore(tempDir(t))
	calls := 0
	fetcher := func(path string) (map[string]string, error) {
		calls++
		return map[string]string{"KEY": "v"}, nil
	}

	// Store an already-expired entry.
	_ = store.Set("secret_app", &cache.Entry{
		Secrets:   map[string]string{"KEY": "old"},
		FetchedAt: time.Now().Add(-10 * time.Minute),
		TTL:       1 * time.Minute,
	})

	secrets, err := cache.Lookup(store, "secret/app", time.Minute, fetcher)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["KEY"] != "v" {
		t.Errorf("expected refreshed value, got %s", secrets["KEY"])
	}
	if calls != 1 {
		t.Errorf("expected 1 fetch call after expiry, got %d", calls)
	}
}

func TestLookupPropagatesFetcherError(t *testing.T) {
	store := cache.NewStore(tempDir(t))
	fetcher := func(path string) (map[string]string, error) {
		return nil, errors.New("vault unavailable")
	}

	_, err := cache.Lookup(store, "secret/app", time.Minute, fetcher)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
