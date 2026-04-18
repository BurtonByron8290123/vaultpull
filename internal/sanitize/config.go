package sanitize

import (
	"os"
	"strings"
)

// FromEnv builds a Policy by reading environment variables.
//
//	VAULTPULL_SANITIZE_STRIP_INVALID  = "true" | "false"  (default true)
//	VAULTPULL_SANITIZE_NORMALIZE_KEYS = "true" | "false"  (default true)
//	VAULTPULL_SANITIZE_TRIM_VALUES    = "true" | "false"  (default true)
//	VAULTPULL_SANITIZE_STRIP_NULLS    = "true" | "false"  (default true)
func FromEnv() Policy {
	p := DefaultPolicy()
	if v := os.Getenv("VAULTPULL_SANITIZE_STRIP_INVALID"); v != "" {
		p.StripInvalidKeys = parseBool(v, true)
	}
	if v := os.Getenv("VAULTPULL_SANITIZE_NORMALIZE_KEYS"); v != "" {
		p.NormalizeKeys = parseBool(v, true)
	}
	if v := os.Getenv("VAULTPULL_SANITIZE_TRIM_VALUES"); v != "" {
		p.TrimValues = parseBool(v, true)
	}
	if v := os.Getenv("VAULTPULL_SANITIZE_STRIP_NULLS"); v != "" {
		p.StripNullBytes = parseBool(v, true)
	}
	return p
}

func parseBool(s string, def bool) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	}
	return def
}
