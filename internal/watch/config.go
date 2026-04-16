package watch

import (
	"errors"
	"time"
)

// Config holds watch-mode settings, typically sourced from the CLI or config file.
type Config struct {
	// Interval between Vault polls.
	Interval time.Duration `mapstructure:"interval" yaml:"interval"`
	// SnapshotPath is where the last-known state is persisted.
	SnapshotPath string `mapstructure:"snapshot_path" yaml:"snapshot_path"`
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:     30 * time.Second,
		SnapshotPath: ".vaultpull-watch.snap",
	}
}

// Validate returns an error when Config contains invalid values.
func (c Config) Validate() error {
	if c.Interval < time.Second {
		return errors.New("watch interval must be at least 1s")
	}
	if c.SnapshotPath == "" {
		return errors.New("watch snapshot_path must not be empty")
	}
	return nil
}
