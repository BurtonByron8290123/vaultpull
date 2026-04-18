// Package quota provides per-path request counting for vaultpull.
//
// It is designed to prevent runaway Vault API usage during a pull run by
// capping the number of fetch attempts per secret path. The limit is
// configurable via the VAULTPULL_QUOTA_MAX_REQUESTS environment variable
// or programmatically through a Policy.
//
// Basic usage:
//
//	tracker, _ := quota.New(quota.DefaultPolicy())
//	if err := tracker.Allow("/secret/data/myapp"); err != nil {
//	    // handle quota exceeded
//	}
package quota
