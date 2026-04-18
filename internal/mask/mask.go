// Package mask provides utilities for masking secret values in output.
package mask

import "strings"

const defaultMask = "********"

// Policy controls how values are masked.
type Policy struct {
	Mask        string
	RevealChars int // number of trailing chars to reveal (0 = none)
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		Mask:        defaultMask,
		RevealChars: 0,
	}
}

// Apply masks the given value according to the policy.
func (p Policy) Apply(value string) string {
	if value == "" {
		return value
	}
	if p.RevealChars <= 0 || p.RevealChars >= len(value) {
		return p.Mask
	}
	suffix := value[len(value)-p.RevealChars:]
	return p.Mask + suffix
}

// ApplyMap masks all values in the provided map, returning a new map.
func (p Policy) ApplyMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = p.Apply(v)
	}
	return out
}

// IsSensitive returns true if the key looks like a secret based on common patterns.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "APIKEY",
	"PRIVATE_KEY", "PRIVATEKEY", "CREDENTIAL", "AUTH",
}
