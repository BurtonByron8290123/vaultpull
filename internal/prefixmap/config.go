package prefixmap

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type fileConfig struct {
	Mappings []Entry `yaml:"mappings"`
}

// LoadConfig reads prefix mappings from a YAML file.
// Expected format:
//
//	mappings:
//	  - path: secret/app/prod
//	    prefix: PROD_
//	  - path: secret/app/shared
//	    prefix: SHARED_
func LoadConfig(path string) (*Mapper, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("prefixmap: open %s: %w", path, err)
	}
	defer f.Close()

	var cfg fileConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("prefixmap: decode %s: %w", path, err)
	}
	return New(cfg.Mappings)
}
