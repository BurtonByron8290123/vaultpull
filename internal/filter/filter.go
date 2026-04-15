package filter

import (
	"strings"
)

// Rule defines a single include or exclude pattern for secret keys.
type Rule struct {
	Prefix  string
	Exclude bool
}

// Filter holds a set of rules to apply when selecting secrets.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of pattern strings.
// Patterns prefixed with '!' are treated as exclusions.
func New(patterns []string) *Filter {
	rules := make([]Rule, 0, len(patterns))
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, "!") {
			rules = append(rules, Rule{Prefix: strings.TrimPrefix(p, "!"), Exclude: true})
		} else {
			rules = append(rules, Rule{Prefix: p, Exclude: false})
		}
	}
	return &Filter{rules: rules}
}

// Apply returns only the key-value pairs from secrets that match the filter.
// If no include rules are defined, all keys are included by default.
// Exclude rules always take precedence over include rules.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	if len(f.rules) == 0 {
		return secrets
	}

	hasIncludes := false
	for _, r := range f.rules {
		if !r.Exclude {
			hasIncludes = true
			break
		}
	}

	result := make(map[string]string)
	for k, v := range secrets {
		if f.excluded(k) {
			continue
		}
		if !hasIncludes || f.included(k) {
			result[k] = v
		}
	}
	return result
}

func (f *Filter) included(key string) bool {
	for _, r := range f.rules {
		if !r.Exclude && strings.HasPrefix(key, r.Prefix) {
			return true
		}
	}
	return false
}

func (f *Filter) excluded(key string) bool {
	for _, r := range f.rules {
		if r.Exclude && strings.HasPrefix(key, r.Prefix) {
			return true
		}
	}
	return false
}
