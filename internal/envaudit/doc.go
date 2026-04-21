// Package envaudit records an append-only audit trail for changes made to
// local .env files by vaultpull.
//
// Each sync operation that modifies secrets appends one JSON line per changed
// key to a configurable log file. Entries capture the timestamp, env file
// path, key name, and the kind of change (added / updated / removed).
//
// Configuration is read from environment variables:
//
//	VAULTPULL_AUDIT_PATH  – destination file (empty = disabled)
//	VAULTPULL_AUDIT_MASK  – whether to mark entries as masked (default true)
//
// The log file is created with 0600 permissions so that secret metadata is
// not world-readable.
package envaudit
