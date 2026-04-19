// Package tagmap rewrites secret keys based on Vault metadata tags.
package tagmap

import (
	"errors"
	"fmt"
	"strings"
)

// Rule maps a Vault metadata tag value to a key prefix replacement.
type Rule struct {
	Tag    string // metadata tag key, e.g. "env"
	Value  string // tag value to match, e.g. "production"
	Prefix string // prefix to prepend to the secret key
}

// Mapper applies tag-based key rewriting rules.
type Mapper struct {
	rules []Rule
}

// New creates a Mapper from the provided rules.
func New(rules []Rule) (*Mapper, error) {
	for i, r := range rules {
		if strings.TrimSpace(r.Tag) == "" {
			return nil, fmt.Errorf("rule %d: tag must not be empty", i)
		}
		if strings.TrimSpace(r.Value) == "" {
			return nil, fmt.Errorf("rule %d: value must not be empty", i)
		}
	}
	return &Mapper{rules: rules}, nil
}

// Apply rewrites key using the first matching rule given the metadata tags.
// If no rule matches, the original key is returned unchanged.
func (m *Mapper) Apply(key string, tags map[string]string) string {
	for _, r := range m.rules {
		if v, ok := tags[r.Tag]; ok && v == r.Value {
			if r.Prefix == "" {
				return key
			}
			return r.Prefix + key
		}
	}
	return key
}

// ApplyMap rewrites all keys in secrets using the provided tags.
func (m *Mapper) ApplyMap(secrets map[string]string, tags map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[m.Apply(k, tags)] = v
	}
	return out
}

// ErrNoRules is returned when no rules are provided.
var ErrNoRules = errors.New("tagmap: no rules provided")
