package envsort

import (
	"fmt"
	"os"
	"strings"
)

const (
	envStrategy  = "VAULTPULL_SORT_STRATEGY"
	envPrefixSep = "VAULTPULL_SORT_PREFIX_SEP"
)

// FromEnv builds a Policy from environment variables, falling back to
// DefaultPolicy() for any value that is absent or unrecognised.
//
//	VAULTPULL_SORT_STRATEGY   – "alpha" | "alpha_desc" | "prefix" (default: "alpha")
//	VAULTPULL_SORT_PREFIX_SEP – separator character for prefix grouping (default: "_")
func FromEnv() (Policy, error) {
	p := DefaultPolicy()

	if raw := strings.TrimSpace(os.Getenv(envStrategy)); raw != "" {
		s, err := parseStrategy(raw)
		if err != nil {
			return p, err
		}
		p.Strategy = s
	}

	if sep := os.Getenv(envPrefixSep); sep != "" {
		p.PrefixSep = sep
	}

	return p, nil
}

func parseStrategy(s string) (Strategy, error) {
	switch strings.ToLower(s) {
	case "alpha", "alphabetical":
		return Alphabetical, nil
	case "alpha_desc", "alphabetical_desc":
		return AlphabeticalDesc, nil
	case "prefix", "prefix_grouped":
		return PrefixGrouped, nil
	default:
		return Alphabetical, fmt.Errorf("envsort: unknown strategy %q", s)
	}
}
