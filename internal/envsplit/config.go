package envsplit

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type rawRule struct {
	Prefix string `yaml:"prefix"`
	Group  string `yaml:"group"`
	Strip  bool   `yaml:"strip"`
}

type rawConfig struct {
	Rules        []rawRule `yaml:"rules"`
	DefaultGroup string    `yaml:"default_group"`
}

// LoadConfig reads a YAML file and returns a Policy.
//
// Example YAML:
//
//	rules:
//	  - prefix: APP_
//	    group: app
//	    strip: true
//	  - prefix: DB_
//	    group: database
//	default_group: misc
func LoadConfig(path string) (Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Policy{}, fmt.Errorf("envsplit: read config %q: %w", path, err)
	}
	var raw rawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return Policy{}, fmt.Errorf("envsplit: parse config %q: %w", path, err)
	}
	p := Policy{DefaultGroup: strings.TrimSpace(raw.DefaultGroup)}
	for _, r := range raw.Rules {
		p.Rules = append(p.Rules, Rule{
			Prefix: r.Prefix,
			Group:  r.Group,
			Strip:  r.Strip,
		})
	}
	return p, nil
}
