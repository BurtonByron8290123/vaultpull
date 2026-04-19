package preflight

import (
	"os"
	"strings"
)

// Config holds the values needed to build the default preflight runner.
type Config struct {
	VaultAddr  string
	VaultToken string
	OutputPath string
}

// FromEnv populates a Config from environment variables, falling back to
// the supplied cfg values when env vars are absent.
func FromEnv(cfg Config) Config {
	if v := strings.TrimSpace(os.Getenv("VAULT_ADDR")); v != "" {
		cfg.VaultAddr = v
	}
	if v := strings.TrimSpace(os.Getenv("VAULT_TOKEN")); v != "" {
		cfg.VaultToken = v
	}
	return cfg
}

// Build constructs a Runner from the Config.
func (c Config) Build() *Runner {
	return Default(c.VaultAddr, c.VaultToken, c.OutputPath)
}
