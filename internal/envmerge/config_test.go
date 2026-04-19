package envmerge

import (
	"testing"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv(envStrategy, "")
	p, err := FromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Strategy != StrategyLastWins {
		t.Errorf("expected LastWins, got %d", p.Strategy)
	}
}

func TestFromEnvReadsStrategy(t *testing.T) {
	cases := []struct {
		val  string
		want Strategy
	}{
		{"last", StrategyLastWins},
		{"last-wins", StrategyLastWins},
		{"first", StrategyFirstWins},
		{"first-wins", StrategyFirstWins},
		{"error", StrategyError},
	}
	for _, tc := range cases {
		t.Run(tc.val, func(t *testing.T) {
			t.Setenv(envStrategy, tc.val)
			p, err := FromEnv()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p.Strategy != tc.want {
				t.Errorf("want %d got %d", tc.want, p.Strategy)
			}
		})
	}
}

func TestFromEnvInvalidStrategyReturnsError(t *testing.T) {
	t.Setenv(envStrategy, "bogus")
	_, err := FromEnv()
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}
