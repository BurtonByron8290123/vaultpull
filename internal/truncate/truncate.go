// Package truncate provides helpers for limiting secret values written to
// .env files — useful when a Vault secret contains a very long token that
// would otherwise bloat the file or exceed shell limits.
package truncate

import "fmt"

// Policy controls how values are truncated.
type Policy struct {
	// MaxLen is the maximum number of runes allowed. 0 means no limit.
	MaxLen int
	// Suffix is appended when a value is truncated (e.g. "…").
	Suffix string
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		MaxLen: 0,
		Suffix: "...",
	}
}

// Apply truncates value according to p. If MaxLen is 0 the original value is
// returned unchanged.
func (p Policy) Apply(value string) string {
	if p.MaxLen <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) <= p.MaxLen {
		return value
	}
	suffix := []rune(p.Suffix)
	cutAt := p.MaxLen - len(suffix)
	if cutAt < 0 {
		cutAt = 0
	}
	return string(runes[:cutAt]) + p.Suffix
}

// ApplyMap truncates every value in m according to p and returns a new map.
func (p Policy) ApplyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = p.Apply(v)
	}
	return out
}

// Validate returns an error when the policy is misconfigured.
func (p Policy) Validate() error {
	if p.MaxLen < 0 {
		return fmt.Errorf("truncate: MaxLen must be >= 0, got %d", p.MaxLen)
	}
	if p.MaxLen > 0 && len([]rune(p.Suffix)) >= p.MaxLen {
		return fmt.Errorf("truncate: Suffix length must be less than MaxLen")
	}
	return nil
}
