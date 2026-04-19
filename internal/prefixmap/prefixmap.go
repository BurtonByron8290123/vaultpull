// Package prefixmap maps Vault secret paths to env key prefixes.
package prefixmap

import (
	"fmt"
	"strings"
)

// Entry maps a Vault path to an optional env key prefix.
type Entry struct {
	Path   string `yaml:"path"`
	Prefix string `yaml:"prefix"`
}

// Mapper rewrites secret keys using configured prefix rules.
type Mapper struct {
	entries []Entry
}

// New returns a Mapper from the given entries.
func New(entries []Entry) (*Mapper, error) {
	for i, e := range entries {
		if strings.TrimSpace(e.Path) == "" {
			return nil, fmt.Errorf("prefixmap: entry %d has empty path", i)
		}
	}
	return &Mapper{entries: entries}, nil
}

// Apply rewrites key using the first entry whose Path matches vaultPath.
// If no entry matches, the key is returned unchanged.
func (m *Mapper) Apply(vaultPath, key string) string {
	for _, e := range m.entries {
		if e.Path == vaultPath && e.Prefix != "" {
			return e.Prefix + key
		}
	}
	return key
}

// ApplyMap rewrites all keys in secrets for the given vaultPath.
func (m *Mapper) ApplyMap(vaultPath string, secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[m.Apply(vaultPath, k)] = v
	}
	return out
}

// Paths returns the list of configured Vault paths in order.
func (m *Mapper) Paths() []string {
	paths := make([]string, len(m.entries))
	for i, e := range m.entries {
		paths[i] = e.Path
	}
	return paths
}
