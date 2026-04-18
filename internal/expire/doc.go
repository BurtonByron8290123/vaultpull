// Package expire provides TTL-based expiry logic for vaultpull secret entries.
//
// Use Policy to determine whether a previously fetched secret should be
// considered stale and re-fetched from Vault.
//
// Configuration can be loaded from environment variables via FromEnv, or
// constructed programmatically using DefaultPolicy.
//
// Example:
//
//	p := expire.DefaultPolicy()
//	if p.IsExpired(entry.FetchedAt) {
//		// re-fetch from Vault
//	}
package expire
