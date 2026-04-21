// Package envtag attaches metadata tags to env keys based on configurable rules.
package envtag

import (
	"fmt"
	"strings"
)

// Tag represents a metadata label attached to an env key.
type Tag struct {
	Key   string
	Value string
}

// Rule maps a key prefix to a set of tags to apply.
type Rule struct {
	Prefix string
	Tags   []Tag
}

// Policy holds the set of tagging rules.
type Policy struct {
	Rules []Rule
}

// Tagger applies tags to env maps.
type Tagger struct {
	policy Policy
}

// New creates a Tagger from the given policy.
func New(p Policy) (*Tagger, error) {
	for i, r := range p.Rules {
		if r.Prefix == "" {
			return nil, fmt.Errorf("envtag: rule[%d] has empty prefix", i)
		}
		for j, t := range r.Tags {
			if t.Key == "" {
				return nil, fmt.Errorf("envtag: rule[%d].tag[%d] has empty key", i, j)
			}
		}
	}
	return &Tagger{policy: p}, nil
}

// Apply returns a new map where each env key that matches a rule prefix has
// additional tag keys injected as "__TAG_<KEY>__<ENVKEY>" entries.
func (t *Tagger) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for envKey := range env {
		for _, rule := range t.policy.Rules {
			if strings.HasPrefix(envKey, rule.Prefix) {
				for _, tag := range rule.Tags {
					tagKey := fmt.Sprintf("__TAG_%s__%s", strings.ToUpper(tag.Key), envKey)
					out[tagKey] = tag.Value
				}
				break
			}
		}
	}
	return out
}

// ApplyMap applies tags only to the provided keys subset.
func (t *Tagger) ApplyMap(env map[string]string, keys []string) map[string]string {
	sub := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := env[k]; ok {
			sub[k] = v
		}
	}
	return t.Apply(sub)
}
