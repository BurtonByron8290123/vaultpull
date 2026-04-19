package envmerge

import (
	"fmt"
	"os"
	"strings"
)

const (
	envStrategy = "VAULTPULL_MERGE_STRATEGY"
)

// FromEnv reads merge configuration from environment variables.
func FromEnv() (Policy, error) {
	p := DefaultPolicy()
	if raw := os.Getenv(envStrategy); raw != "" {
		s, err := parseStrategy(raw)
		if err != nil {
			return p, err
		}
		p.Strategy = s
	}
	return p, nil
}

func parseStrategy(s string) (Strategy, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "last", "last-wins":
		return StrategyLastWins, nil
	case "first", "first-wins":
		return StrategyFirstWins, nil
	case "error":
		return StrategyError, nil
	default:
		return 0, fmt.Errorf("envmerge: unknown strategy %q (want last|first|error)", s)
	}
}
