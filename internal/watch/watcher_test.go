package watch_test

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/vaultpull/internal/snapshot"
	"github.com/yourorg/vaultpull/internal/watch"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "watch-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestWatcherDetectsChange(t *testing.T) {
	dir := tempDir(t)
	snap := snapshot.NewStore()
	snapshotPath := filepath.Join(dir, "snap.json")

	calls := int32(0)
	fetcher := func(ctx context.Context) (map[string]string, error) {
		return map[string]string{"KEY": "value"}, nil
	}
	onChange := func(_, _ map[string]string) { atomic.AddInt32(&calls, 1) }

	w := watch.New(20*time.Millisecond, snapshotPath, snap, fetcher, onChange)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	_ = w.Run(ctx)

	if atomic.LoadInt32(&calls) == 0 {
		t.Fatal("expected onChange to be called at least once")
	}
}

func TestWatcherNoCallbackWhenUnchanged(t *testing.T) {
	dir := tempDir(t)
	snap := snapshot.NewStore()
	snapshotPath := filepath.Join(dir, "snap.json")

	secrets := map[string]string{"KEY": "value"}
	// Pre-populate snapshot so first poll sees no change.
	if err := snap.Save(snapshotPath, snapshot.Build(secrets)); err != nil {
		t.Fatal(err)
	}

	calls := int32(0)
	fetcher := func(ctx context.Context) (map[string]string, error) { return secrets, nil }
	onChange := func(_, _ map[string]string) { atomic.AddInt32(&calls, 1) }

	w := watch.New(20*time.Millisecond, snapshotPath, snap, fetcher, onChange)
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	_ = w.Run(ctx)

	if atomic.LoadInt32(&calls) != 0 {
		t.Fatalf("expected no onChange calls, got %d", calls)
	}
}

func TestWatcherCancelStopsLoop(t *testing.T) {
	dir := tempDir(t)
	snap := snapshot.NewStore()
	snapshotPath := filepath.Join(dir, "snap.json")

	fetcher := func(ctx context.Context) (map[string]string, error) { return map[string]string{}, nil }
	w := watch.New(10*time.Millisecond, snapshotPath, snap, fetcher, func(_, _ map[string]string) {})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := w.Run(ctx)
	if err == nil {
		t.Fatal("expected context error")
	}
}
