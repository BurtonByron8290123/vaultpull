package config

import (
	"fmt"
	"time"

	"github.com/user/vaultpull/internal/hooks"
)

// HooksConfig mirrors hooks.Config for YAML/mapstructure unmarshalling.
type HooksConfig struct {
	PrePull  string        `mapstructure:"pre_pull"  yaml:"pre_pull"`
	PostPull string        `mapstructure:"post_pull" yaml:"post_pull"`
	Timeout  time.Duration `mapstructure:"timeout"   yaml:"timeout"`
}

// ToHooksConfig converts the embedded HooksConfig into a hooks.Config.
func (c *Config) ToHooksConfig() hooks.Config {
	return hooks.Config{
		PrePull:  c.Hooks.PrePull,
		PostPull: c.Hooks.PostPull,
		Timeout:  c.Hooks.Timeout,
	}
}

// validateHooks checks that hook commands are non-empty strings when set and
// that the timeout, if provided, is a positive duration.
func validateHooks(h HooksConfig) error {
	if h.Timeout < 0 {
		return fmt.Errorf("config: hooks.timeout must be a positive duration, got %s", h.Timeout)
	}
	return nil
}
