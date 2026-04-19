package depcheck

import (
	"os"
	"strings"
)

const envVar = "VAULTPULL_REQUIRED_KEYS"

// FromEnv builds a Policy by reading a comma-separated list of required keys
// from the VAULTPULL_REQUIRED_KEYS environment variable.
//
// Example:
//
//	export VAULTPULL_REQUIRED_KEYS=DATABASE_URL,API_KEY,SECRET_TOKEN
func FromEnv() Policy {
	raw := strings.TrimSpace(os.Getenv(envVar))
	if raw == "" {
		return Policy{}
	}
	parts := strings.Split(raw, ",")
	keys := make([]string, 0, len(parts))
	for _, p := range parts {
		k := strings.TrimSpace(p)
		if k != "" {
			keys = append(keys, k)
		}
	}
	return Policy{Required: keys}
}
