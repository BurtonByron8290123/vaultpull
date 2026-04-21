package envcheck

import (
	"os"
	"strings"
)

const (
	envRequired = "VAULTPULL_CHECK_REQUIRED"
	envNonEmpty = "VAULTPULL_CHECK_NONEMPTY"
)

// FromEnv builds a Policy from environment variables.
//
//	VAULTPULL_CHECK_REQUIRED  – comma-separated list of keys that must exist.
//	VAULTPULL_CHECK_NONEMPTY  – comma-separated list of keys that must be non-empty.
func FromEnv() Policy {
	return Policy{
		Required: splitCSV(os.Getenv(envRequired)),
		NonEmpty: splitCSV(os.Getenv(envNonEmpty)),
	}
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
