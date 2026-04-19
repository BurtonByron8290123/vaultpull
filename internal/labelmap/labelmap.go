// Package labelmap rewrites secret keys using a static label mapping.
// If a key matches an entry in the map, it is renamed to the configured label.
package labelmap

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Entry maps an original key to a target label.
type Entry struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

// Mapper holds the ordered list of rename rules.
type Mapper struct {
	rules []Entry
}

// New creates a Mapper from a slice of entries.
func New(rules []Entry) (*Mapper, error) {
	for i, r := range rules {
		if r.From == "" {
			return nil, fmt.Errorf("labelmap: entry %d has empty 'from'", i)
		}
		if r.To == "" {
			return nil, fmt.Errorf("labelmap: entry %d has empty 'to'", i)
		}
	}
	return &Mapper{rules: rules}, nil
}

// Apply renames keys in src according to the rules and returns a new map.
// Keys without a matching rule are passed through unchanged.
func (m *Mapper) Apply(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	for _, r := range m.rules {
		if v, ok := out[r.From]; ok {
			delete(out, r.From)
			out[r.To] = v
		}
	}
	return out
}

// LoadConfig reads a YAML file and returns a Mapper.
func LoadConfig(path string) (*Mapper, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("labelmap: read config: %w", err)
	}
	var rules []Entry
	if err := yaml.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("labelmap: parse config: %w", err)
	}
	return New(rules)
}
