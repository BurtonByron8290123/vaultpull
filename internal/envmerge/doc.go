// Package envmerge provides a configurable strategy for merging
// multiple env maps fetched from different Vault secret paths.
//
// Three strategies are supported:
//
//   - last-wins  (default): the last source defining a key wins.
//   - first-wins: the first source defining a key wins.
//   - error: any key conflict with differing values returns an error.
//
// The active strategy can be set via the VAULTPULL_MERGE_STRATEGY
// environment variable (values: last, first, error).
package envmerge
