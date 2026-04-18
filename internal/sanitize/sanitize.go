// Package sanitize provides key and value sanitization for env variables
// before they are written to disk.
package sanitize

import (
	"regexp"
	"strings"
)

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Policy controls sanitization behaviour.
type Policy struct {
	// StripInvalidKeys drops keys that are not valid shell identifiers.
	StripInvalidKeys bool
	// NormalizeKeys uppercases all keys.
	NormalizeKeys bool
	// TrimValues trims leading/trailing whitespace from values.
	TrimValues bool
	// StripNullBytes removes null bytes from values.
	StripNullBytes bool
}

// DefaultPolicy returns a Policy with safe defaults.
func DefaultPolicy() Policy {
	return Policy{
		StripInvalidKeys: true,
		NormalizeKeys:    true,
		TrimValues:       true,
		StripNullBytes:   true,
	}
}

// Sanitizer applies a Policy to env maps.
type Sanitizer struct {
	policy Policy
}

// New creates a Sanitizer with the given policy.
func New(p Policy) *Sanitizer {
	return &Sanitizer{policy: p}
}

// Apply returns a sanitized copy of the input map.
func (s *Sanitizer) Apply(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		if s.policy.NormalizeKeys {
			k = strings.ToUpper(k)
		}
		if s.policy.StripInvalidKeys && !validKeyRe.MatchString(k) {
			continue
		}
		if s.policy.TrimValues {
			v = strings.TrimSpace(v)
		}
		if s.policy.StripNullBytes {
			v = strings.ReplaceAll(v, "\x00", "")
		}
		out[k] = v
	}
	return out
}

// IsValidKey reports whether k is a valid shell environment variable name.
func IsValidKey(k string) bool {
	return validKeyRe.MatchString(k)
}
