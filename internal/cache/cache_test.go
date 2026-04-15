package cache_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/cache"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestSetAndGet(t *testing.T) {
	store := cache.NewStore(tempDir(t))
	entry := &cache.Entry{
		Secrets:   map[string]string{"FOO": "bar"},
		FetchedAt: time.Now(),
		TTL:       5 * time.Minute,
	}
	if err := store.Set("mypath", entry); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := store.Get("mypath")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil {
		t.Fatal("expected entry, got nil")
	}
	if got.Secrets["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", got.Secrets["FOO"])
	}
}

func TestGetMissingReturnsNil(t *testing.T) {
	store := cache.NewStore(tempDir(t))
	got, err := store.Get("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestInvalidateRemovesEntry(t *testing.T) {
	store := cache.NewStore(tempDir(t))
	_ = store.Set("k", &cache.Entry{Secrets: map[string]string{"A": "1"}, FetchedAt: time.Now()})
	if err := store.Invalidate("k"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	got, _ := store.Get("k")
	if got != nil {
		t.Error("expected nil after invalidation")
	}
}

func TestInvalidateNoopWhenMissing(t *testing.T) {
	store := cache.NewStore(tempDir(t))
	if err := store.Invalidate("ghost"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestIsExpired(t *testing.T) {
	expired := &cache.Entry{FetchedAt: time.Now().Add(-10 * time.Minute), TTL: 5 * time.Minute}
	if !expired.IsExpired() {
		t.Error("expected entry to be expired")
	}
	fresh := &cache.Entry{FetchedAt: time.Now(), TTL: 5 * time.Minute}
	if fresh.IsExpired() {
		t.Error("expected entry to be fresh")
	}
}

func TestZeroTTLNeverExpires(t *testing.T) {
	e := &cache.Entry{FetchedAt: time.Now().Add(-24 * time.Hour), TTL: 0}
	if e.IsExpired() {
		t.Error("zero TTL should never expire")
	}
}

func TestFilePermissions(t *testing.T) {
	dir := tempDir(t)
	store := cache.NewStore(dir)
	_ = store.Set("perms", &cache.Entry{Secrets: map[string]string{}, FetchedAt: time.Now()})

	info, err := os.Stat(dir + "/perms.json")
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected 0600, got %o", perm)
	}
}
