package labelmap

import (
	"fmt"
	"os"
	"strings"
)

const envKey = "VAULTPULL_LABELMAP_FILE"

// FromEnv loads a Mapper from the path specified in VAULTPULL_LABELMAP_FILE.
// If the variable is unset or empty, a no-op Mapper is returned.
// Returns an error if the variable is set but the file cannot be loaded.
func FromEnv() (*Mapper, error) {
	path := strings.TrimSpace(os.Getenv(envKey))
	if path == "" {
		return &Mapper{}, nil
	}
	m, err := LoadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("labelmap: FromEnv: %s=%q: %w", envKey, path, err)
	}
	return m, nil
}

// EnvPath returns the label map file path configured via the environment
// variable VAULTPULL_LABELMAP_FILE, trimmed of surrounding whitespace.
// Returns an empty string if the variable is unset or blank.
func EnvPath() string {
	return strings.TrimSpace(os.Getenv(envKey))
}
