package envexport

import (
	"fmt"
	"os"
	"strings"
)

const envVarFormat = "VAULTPULL_EXPORT_FORMAT"

// Config holds envexport settings.
type Config struct {
	Format Format
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{Format: FormatDotenv}
}

// FromEnv builds a Config from environment variables, falling back to
// defaults for missing or invalid values.
func FromEnv() (Config, error) {
	cfg := DefaultConfig()
	if raw := strings.TrimSpace(os.Getenv(envVarFormat)); raw != "" {
		f := Format(strings.ToLower(raw))
		switch f {
		case FormatDotenv, FormatJSON, FormatExport:
			cfg.Format = f
		default:
			return Config{}, fmt.Errorf("envexport: invalid format in %s: %q", envVarFormat, raw)
		}
	}
	return cfg, nil
}
