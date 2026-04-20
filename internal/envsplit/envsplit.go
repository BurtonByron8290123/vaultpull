// Package envsplit splits a flat env map into multiple named groups
// based on prefix rules defined in a configuration file or inline config.
package envsplit

import (
	"fmt"
	"strings"
)

// Rule maps a key prefix to a named group.
type Rule struct {
	Prefix string
	Group  string
	Strip  bool // if true, remove the prefix from the key in the output
}

// Policy holds the set of split rules.
type Policy struct {
	Rules       []Rule
	DefaultGroup string // keys that match no rule go here; ignored if empty
}

// Splitter partitions an env map according to a Policy.
type Splitter struct {
	policy Policy
}

// New returns a Splitter for the given Policy.
func New(p Policy) (*Splitter, error) {
	for i, r := range p.Rules {
		if r.Prefix == "" {
			return nil, fmt.Errorf("envsplit: rule %d has empty prefix", i)
		}
		if r.Group == "" {
			return nil, fmt.Errorf("envsplit: rule %d has empty group", i)
		}
	}
	return &Splitter{policy: p}, nil
}

// Apply partitions src into groups. Each returned map contains only the
// keys that matched the corresponding rule. Rules are evaluated in order;
// the first matching rule wins.
func (s *Splitter) Apply(src map[string]string) map[string]map[string]string {
	out := make(map[string]map[string]string)

	for k, v := range src {
		matched := false
		for _, r := range s.policy.Rules {
			if strings.HasPrefix(k, r.Prefix) {
				if out[r.Group] == nil {
					out[r.Group] = make(map[string]string)
				}
				outKey := k
				if r.Strip {
					outKey = strings.TrimPrefix(k, r.Prefix)
				}
				out[r.Group][outKey] = v
				matched = true
				break
			}
		}
		if !matched && s.policy.DefaultGroup != "" {
			if out[s.policy.DefaultGroup] == nil {
				out[s.policy.DefaultGroup] = make(map[string]string)
			}
			out[s.policy.DefaultGroup][k] = v
		}
	}
	return out
}
