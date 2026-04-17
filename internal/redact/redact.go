// Package redact provides utilities for masking sensitive values
// before they are printed to logs or terminal output.
package redact

import "strings"

const defaultMask = "****"

// Redactor masks secret values based on a set of sensitive keys.
type Redactor struct {
	keys map[string]struct{}
	mask string
}

// New returns a Redactor that will mask values whose keys match any
// of the provided sensitive key names (case-insensitive).
func New(sensitiveKeys []string) *Redactor {
	km := make(map[string]struct{}, len(sensitiveKeys))
	for _, k := range sensitiveKeys {
		km[strings.ToUpper(k)] = struct{}{}
	}
	return &Redactor{keys: km, mask: defaultMask}
}

// WithMask returns a copy of the Redactor using a custom mask string.
func (r * string) *Redactor {
	return &Redactor{keys: r.keys, mask: mask}
}

// Value returns the masked value if the key is sensitive, otherwise the
// original value.
func (r *Redactor) Value(key, value string) string {
	if r.IsSensitive(key) {
		return r.mask
	}
	return value
}

// IsSensitive reports whether the given key should be treated as sensitive.
func (r *Redactor) IsSensitive(key string) bool {
	_, ok := r.keys[strings.ToUpper(key)]
	return ok
}

// Map returns a copy of m with sensitive values replaced by the mask.
func (r *Redactor) Map(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = r.Value(k, v)
	}
	return out
}
