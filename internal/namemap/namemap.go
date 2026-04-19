// Package namemap rewrites secret keys according to a static mapping.
package namemap

import (
	"encoding/json"
	"fmt"
	"os"
)

// Rule maps one key name to another.
type Rule struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Mapper holds a set of rename rules.
type Mapper struct {
	rules []Rule
}

// New creates a Mapper from the given rules. Returns an error if any rule has
// a blank From or To field.
func New(rules []Rule) (*Mapper, error) {
	for i, r := range rules {
		if r.From == "" {
			return nil, fmt.Errorf("namemap: rule %d has empty from", i)
		}
		if r.To == "" {
			return nil, fmt.Errorf("namemap: rule %d has empty to", i)
		}
	}
	return &Mapper{rules: rules}, nil
}

// LoadConfig reads a JSON file containing an array of Rule objects.
func LoadConfig(path string) (*Mapper, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("namemap: read config: %w", err)
	}
	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("namemap: parse config: %w", err)
	}
	return New(rules)
}

// Apply renames the key according to the first matching rule. If no rule
// matches, the original key is returned unchanged.
func (m *Mapper) Apply(key string) string {
	for _, r := range m.rules {
		if r.From == key {
			return r.To
		}
	}
	return key
}

// ApplyMap rewrites all keys in src that have a matching rule and returns a
// new map. Values are preserved; later rules do not chain.
func (m *Mapper) ApplyMap(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[m.Apply(k)] = v
	}
	return out
}
