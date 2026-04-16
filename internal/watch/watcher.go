// Package watch polls Vault at a fixed interval and triggers a callback
// when secrets change relative to the last known snapshot.
package watch

import (
	"context"
	"time"

	"github.com/yourorg/vaultpull/internal/snapshot"
)

// FetchFunc retrieves the current set of secrets from Vault.
type FetchFunc func(ctx context.Context) (map[string]string, error)

// ChangeFunc is called whenever a change is detected.
type ChangeFunc func(prev, next map[string]string)

// Watcher polls Vault on a fixed interval.
type Watcher struct {
	interval time.Duration
	fetch    FetchFunc
	onChange ChangeFunc
	store    *snapshot.Store
	path     string
}

// New creates a Watcher. path is used to persist the snapshot between runs.
func New(interval time.Duration, path string, store *snapshot.Store, fetch FetchFunc, onChange ChangeFunc) *Watcher {
	return &Watcher{
		interval: interval,
		fetch:    fetch,
		onChange: onChange,
		store:    store,
		path:     path,
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.poll(ctx); err != nil {
				return err
			}
		}
	}
}

func (w *Watcher) poll(ctx context.Context) error {
	secrets, err := w.fetch(ctx)
	if err != nil {
		return err
	}
	prev, err := w.store.Load(w.path)
	if err != nil {
		return err
	}
	var prevMap map[string]string
	if prev != nil {
		prevMap = prev.ToMap()
	}
	changes := snapshot.Diff(prev, snapshot.Build(secrets))
	if len(changes) > 0 {
		w.onChange(prevMap, secrets)
	}
	return w.store.Save(w.path, snapshot.Build(secrets))
}
