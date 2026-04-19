package envpin

import (
	"os"
	"path/filepath"
)

const (
	defaultPinFile = ".vaultpull.pins.json"
	envPinFile     = "VAULTPULL_PIN_FILE"
	envPinEnabled  = "VAULTPULL_PIN_ENABLED"
)

// Config holds runtime settings for the pin store.
type Config struct {
	// FilePath is the path to the JSON pin file.
	FilePath string
	// Enabled controls whether pin tracking is active.
	Enabled bool
}

// DefaultConfig returns a Config populated with defaults.
func DefaultConfig() Config {
	return Config{
		FilePath: defaultPinFile,
		Enabled:  true,
	}
}

// FromEnv builds a Config from environment variables, falling back to defaults.
func FromEnv() Config {
	cfg := DefaultConfig()
	if v := os.Getenv(envPinFile); v != "" {
		cfg.FilePath = v
	}
	if v := os.Getenv(envPinEnabled); v == "false" || v == "0" {
		cfg.Enabled = false
	}
	return cfg
}

// Resolve returns an absolute path for cfg.FilePath relative to baseDir when
// the path is not already absolute.
func (c Config) Resolve(baseDir string) string {
	if filepath.IsAbs(c.FilePath) {
		return c.FilePath
	}
	return filepath.Join(baseDir, c.FilePath)
}
