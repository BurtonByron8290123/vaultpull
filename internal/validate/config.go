package validate

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type ruleYAML struct {
	Key      string   `yaml:"key"`
	MinLen   int      `yaml:"min_len"`
	Disallow []string `yaml:"disallow"`
}

type policyYAML struct {
	Rules []ruleYAML `yaml:"rules"`
}

// FromFile loads a Policy from a YAML file at path.
func FromFile(path string) (Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Policy{}, fmt.Errorf("validate: read config: %w", err)
	}
	var py policyYAML
	if err := yaml.Unmarshal(data, &py); err != nil {
		return Policy{}, fmt.Errorf("validate: parse config: %w", err)
	}
	p := Policy{}
	for _, r := range py.Rules {
		if r.Key == "" {
			return Policy{}, fmt.Errorf("validate: rule missing key")
		}
		p.Rules = append(p.Rules, Rule{Key: r.Key, MinLen: r.MinLen, Disallow: r.Disallow})
	}
	return p, nil
}

// FromEnv builds a Policy from VAULTPULL_VALIDATE_REQUIRE (comma-separated
// keys) and VAULTPULL_VALIDATE_MIN_LEN (integer applied to all required keys).
func FromEnv() Policy {
	p := Policy{}
	require := os.Getenv("VAULTPULL_VALIDATE_REQUIRE")
	if require == "" {
		return p
	}
	minLen := 0
	if ml := os.Getenv("VAULTPULL_VALIDATE_MIN_LEN"); ml != "" {
		if n, err := strconv.Atoi(ml); err == nil {
			minLen = n
		}
	}
	for _, k := range strings.Split(require, ",") {
		k = strings.TrimSpace(k)
		if k != "" {
			p.Rules = append(p.Rules, Rule{Key: k, MinLen: minLen})
		}
	}
	return p
}
