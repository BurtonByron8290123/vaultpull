package envcompare

import (
	"os"
	"strconv"
)

// Config controls envcompare behaviour.
type Config struct {
	// MaskValues hides secret values in output.
	MaskValues bool
	// MaskMark is the string used to replace masked values.
	MaskMark string
}

// DefaultConfig returns safe defaults.
func DefaultConfig() Config {
	return Config{
		MaskValues: true,
		MaskMark:   "***",
	}
}

// FromEnv reads configuration from environment variables.
//
//	VAULTPULL_COMPARE_MASK_VALUES  – "true"/"false" (default true)
//	VAULTPULL_COMPARE_MASK_MARK    – replacement string (default "***")
func FromEnv() Config {
	cfg := DefaultConfig()
	if v := os.Getenv("VAULTPULL_COMPARE_MASK_VALUES"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.MaskValues = b
		}
	}
	if v := os.Getenv("VAULTPULL_COMPARE_MASK_MARK"); v != "" {
		cfg.MaskMark = v
	}
	return cfg
}

// Apply returns a copy of Result with values masked according to cfg.
func (cfg Config) Apply(r *Result) *Result {
	if !cfg.MaskValues {
		return r
	}
	masked := make([]Change, len(r.Changes))
	for i, c := range r.Changes {
		c.OldValue = cfg.MaskMark
		c.NewValue = cfg.MaskMark
		masked[i] = c
	}
	return &Result{Changes: masked}
}
