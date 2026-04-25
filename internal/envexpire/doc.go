// Package envexpire provides age-based expiry checking for secrets stored
// in .env files.
//
// A Checker evaluates a map of key → last-rotated timestamps against a
// configurable Policy. Keys are classified as OK, Warning (approaching
// expiry), or Expired so that callers can decide whether to trigger a
// rotation or merely warn the operator.
//
// Configuration can be loaded from environment variables via FromEnv, or
// constructed manually with DefaultPolicy.
package envexpire
