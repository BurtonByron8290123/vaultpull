package envrewrite

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig reads a YAML file containing rewrite rules.
//
// Example YAML:
//
//	rules:
//	  - pattern: "^http://"
//	    replacement: "https://"
//	  - pattern: "localhost"
//	    replacement: "prod.internal"
//	    key_glob: "DATABASE_"
func LoadConfig(path string) (Policy, error) {
	f, err := os.Open(path)
	if err != nil {
		return Policy{}, fmt.Errorf("envrewrite: open config %q: %w", path, err)
	}
	defer f.Close()

	var p Policy
	if err := yaml.NewDecoder(f).Decode(&p); err != nil {
		return Policy{}, fmt.Errorf("envrewrite: decode config %q: %w", path, err)
	}
	return p, nil
}

// FromEnv builds a Rewriter from the YAML file referenced by the
// VAULTPULL_REWRITE_CONFIG environment variable. Returns a no-op
// Rewriter when the variable is unset.
func FromEnv() (*Rewriter, error) {
	path := os.Getenv("VAULTPULL_REWRITE_CONFIG")
	if path == "" {
		return New(Policy{})
	}
	p, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}
	return New(p)
}
