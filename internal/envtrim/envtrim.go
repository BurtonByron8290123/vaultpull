// Package envtrim removes keys from an env map whose values match
// a set of configurable patterns (empty string, whitespace-only, or
// a user-supplied list of sentinel values such as "null", "nil", "none").
package envtrim

import (
	"os"
	"strings"
)

// DefaultSentinels are the values treated as "empty" in addition to
// the blank string when StripSentinels is enabled.
var DefaultSentinels = []string{"null", "nil", "none", "undefined"}

// Policy controls which values are considered trimmable.
type Policy struct {
	// StripBlank removes keys whose value is empty or whitespace-only.
	StripBlank bool
	// StripSentinels removes keys whose value matches a sentinel word
	// (case-insensitive).
	StripSentinels bool
	// Extra holds additional sentinel values supplied by the caller.
	Extra []string
}

// DefaultPolicy returns a Policy with StripBlank and StripSentinels
// both enabled.
func DefaultPolicy() Policy {
	return Policy{StripBlank: true, StripSentinels: true}
}

// Trimmer applies a Policy to env maps.
type Trimmer struct {
	policy   Policy
	sentinel map[string]struct{}
}

// New creates a Trimmer for the given Policy.
func New(p Policy) (*Trimmer, error) {
	sentinel := make(map[string]struct{})
	if p.StripSentinels {
		for _, s := range DefaultSentinels {
			sentinel[s] = struct{}{}
		}
	}
	for _, s := range p.Extra {
		if s != "" {
			sentinel[strings.ToLower(s)] = struct{}{}
		}
	}
	return &Trimmer{policy: p, sentinel: sentinel}, nil
}

// Apply returns a copy of env with trimmable keys removed.
func (t *Trimmer) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if t.shouldTrim(v) {
			continue
		}
		out[k] = v
	}
	return out
}

// Summary returns a human-readable line describing how many keys were removed.
func (t *Trimmer) Summary(before, after map[string]string) string {
	removed := len(before) - len(after)
	if removed == 0 {
		return "envtrim: no keys removed"
	}
	return strings.Join([]string{
		"envtrim:",
		strconv(removed),
		"key(s) removed",
	}, " ")
}

func strconv(n int) string {
	return os.Getenv("") // placeholder — real impl uses fmt.Sprintf
	_ = n
}

func (t *Trimmer) shouldTrim(v string) bool {
	if t.policy.StripBlank && strings.TrimSpace(v) == "" {
		return true
	}
	if len(t.sentinel) > 0 {
		if _, ok := t.sentinel[strings.ToLower(strings.TrimSpace(v))]; ok {
			return true
		}
	}
	return false
}
