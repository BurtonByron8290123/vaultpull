package envstamp

import (
	"os"
	"strconv"
)

const (
	envEnabled = "VAULTPULL_STAMP_ENABLED"
	envVersion = "VAULTPULL_STAMP_VERSION_LABEL"
	envSource  = "VAULTPULL_STAMP_SOURCE_PATH"
)

// FromEnv builds a Policy from environment variables, falling back to
// DefaultPolicy() for any value that is absent or unparseable.
//
//	VAULTPULL_STAMP_ENABLED       – "true"/"false" (default: true)
//	VAULTPULL_STAMP_VERSION_LABEL – arbitrary string  (default: "1")
//	VAULTPULL_STAMP_SOURCE_PATH   – vault secret path (default: "")
func FromEnv() Policy {
	p := DefaultPolicy()

	if raw := os.Getenv(envEnabled); raw != "" {
		if b, err := strconv.ParseBool(raw); err == nil {
			p.Enabled = b
		}
	}

	if v := os.Getenv(envVersion); v != "" {
		p.Version = v
	}

	if s := os.Getenv(envSource); s != "" {
		p.Source = s
	}

	return p
}
