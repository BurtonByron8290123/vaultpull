package envtag

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const envVarConfig = "VAULTPULL_ENVTAG_CONFIG"

type rawRule struct {
	Prefix string            `json:"prefix"`
	Tags   map[string]string `json:"tags"`
}

// FromEnv loads a Policy from the path stored in VAULTPULL_ENVTAG_CONFIG.
// If the env var is unset an empty Policy is returned.
func FromEnv() (Policy, error) {
	path := strings.TrimSpace(os.Getenv(envVarConfig))
	if path == "" {
		return Policy{}, nil
	}
	return FromFile(path)
}

// FromFile parses a JSON file containing an array of tagging rules.
func FromFile(path string) (Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Policy{}, fmt.Errorf("envtag: read config %q: %w", path, err)
	}
	var raw []rawRule
	if err := json.Unmarshal(data, &raw); err != nil {
		return Policy{}, fmt.Errorf("envtag: parse config %q: %w", path, err)
	}
	var rules []Rule
	for _, r := range raw {
		rule := Rule{Prefix: r.Prefix}
		for k, v := range r.Tags {
			rule.Tags = append(rule.Tags, Tag{Key: k, Value: v})
		}
		rules = append(rules, rule)
	}
	return Policy{Rules: rules}, nil
}
