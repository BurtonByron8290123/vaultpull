package labelmap

import (
	"fmt"
	"os"
	"strings"
)

const envKey = "VAULTPULL_LABELMAP_FILE"

// FromEnv loads a Mapper from the path specified in VAULTPULL_LABELMAP_FILE.
// If the variable is unset or empty, a no-op Mapper is returned.
func FromEnv() (*Mapper, error) {
	path := strings.TrimSpace(os.Getenv(envKey))
	if path == "" {
		return &Mapper{}, nil
	}
	m, err := LoadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("labelmap: FromEnv: %w", err)
	}
	return m, nil
}
