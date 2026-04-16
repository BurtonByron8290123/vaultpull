// Package notify provides a lightweight structured notification layer for
// vaultpull CLI output. It supports info, warn, and error severity levels
// and can be silenced via a quiet flag for use in scripts or CI pipelines.
//
// Usage:
//
//	n := notify.New("pull", false)
//	n.Infof("synced %d secrets to %s", count, path)
//	n.Warn("token expiry is within 24 hours")
package notify
