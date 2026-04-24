// Package envresolve expands shell-style variable references (${VAR} or $VAR)
// within the values of an env map.
//
// Resolution order:
//  1. Caller-supplied overrides.
//  2. Other keys present in the same env map.
//  3. Process environment (when AllowEnvFallback is true).
//
// Unresolvable placeholders are left in place unless ErrorOnMissing is set,
// in which case Apply returns an error on the first missing reference.
package envresolve
