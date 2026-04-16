// Package watch implements a polling-based watch mode for vaultpull.
//
// It periodically fetches secrets from Vault, compares them against the last
// known snapshot, and invokes a user-supplied callback whenever a change is
// detected. The snapshot is persisted to disk so that restarts do not trigger
// spurious change notifications.
//
// Typical usage:
//
//	w := watch.New(
//		cfg.Interval,
//		cfg.SnapshotPath,
//		snapshotStore,
//		vaultFetcher,
//		func(prev, next map[string]string) { /* re-write .env */ },
//	)
//	if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package watch
