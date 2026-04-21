// Package envrewrite applies regex-based rewrite rules to env values.
package envrewrite

import (
	"fmt"
	"regexp"
)

// Rule describes a single rewrite operation.
type Rule struct {
	Pattern     string `yaml:"pattern"`
	Replacement string `yaml:"replacement"`
	KeyGlob     string `yaml:"key_glob"` // optional: only apply to matching keys

	re *regexp.Regexp
}

// Policy holds a set of rewrite rules.
type Policy struct {
	Rules []Rule
}

// Rewriter applies rewrite rules to env maps.
type Rewriter struct {
	policy Policy
}

// New compiles all rules in p and returns a Rewriter.
func New(p Policy) (*Rewriter, error) {
	for i, r := range p.Rules {
		if r.Pattern == "" {
			return nil, fmt.Errorf("rule %d: pattern must not be empty", i)
		}
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, fmt.Errorf("rule %d: invalid pattern %q: %w", i, r.Pattern, err)
		}
		p.Rules[i].re = re
	}
	return &Rewriter{policy: p}, nil
}

// Apply rewrites values in src according to the configured rules.
// It returns a new map; src is not mutated.
func (rw *Rewriter) Apply(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	for _, rule := range rw.policy.Rules {
		for k, v := range out {
			if !matchesGlob(rule.KeyGlob, k) {
				continue
			}
			out[k] = rule.re.ReplaceAllString(v, rule.Replacement)
		}
	}
	return out
}

// matchesGlob returns true when glob is empty or the key has the glob as prefix.
func matchesGlob(glob, key string) bool {
	if glob == "" {
		return true
	}
	matched, err := regexp.MatchString("^"+regexp.QuoteMeta(glob), key)
	return err == nil && matched
}
