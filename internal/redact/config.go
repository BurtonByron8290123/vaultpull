package redact

import (
	"os"
	"strings"
)

// DefaultSensitiveKeys is the built-in list of key patterns considered sensitive.
var DefaultSensitiveKeys = []string{
	"PASSWORD",
	"PASSWD",
	"SECRET",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIALS",
	"ACCESS",
	"DSN",
}

// Config holds redaction configuration.
type Config struct {
	// ExtraKeys are additional key names to treat as sensitive.
	ExtraKeys []string `yaml:"extra_keys"`
	// Mask is the string used to replace sensitive values.
	Mask string `yaml:"mask"`
}

// FromConfig builds a Redactor from the provided Config, merging
// DefaultSensitiveKeys with any extra keys defined by the user.
func FromConfig(cfg Config) *Redactor {
	keys := make([]string, len(DefaultSensitiveKeys))
	copy(keys, DefaultSensitiveKeys)
	keys = append(keys, cfg.ExtraKeys...)
	r := New(keys)
	if cfg.Mask != "" {
		r = r.WithMask(cfg.Mask)
	}
	return r
}

// FromEnv readsAULTPULL_REDACT_KEYS from the environment (comma-separated)
// and merges them with DefaultSensitiveKeys.
func FromEnv() *Redactor {
	extra := os.Getenv("VAULTPULL_REDACT_KEYS")
	var extras []string
	if extra != "" {
		for _, k := range strings.Split(extra, ",") {
			if k = strings.TrimSpace(k); k != "" {
				extras = append(extras, k)
			}
		}
	}
	return FromConfig(Config{ExtraKeys: extras})
}
