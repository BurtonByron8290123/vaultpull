package multienv

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the multi-env configuration loaded from a YAML file.
type Config struct {
	Targets []Target `yaml:"targets"`
}

// LoadConfig reads a YAML file at path and returns the parsed Config.
// The file is expected to contain a top-level "targets" list.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("multienv: read config %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("multienv: parse config %q: %w", path, err)
	}
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// validateConfig checks that each target has a non-empty name and path.
func validateConfig(cfg *Config) error {
	for i, t := range cfg.Targets {
		if t.Name == "" {
			return fmt.Errorf("multienv: target[%d] missing name", i)
		}
		if t.Path == "" {
			return fmt.Errorf("multienv: target[%d] (%q) missing path", i, t.Name)
		}
	}
	return nil
}
