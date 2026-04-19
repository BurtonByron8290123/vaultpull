package preflight_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/preflight"
)

func TestFromEnvUsesSuppliedValuesWhenEnvAbsent(t *testing.T) {
	cfg := preflight.FromEnv(preflight.Config{
		VaultAddr:  "http://localhost:8200",
		VaultToken: "s.test",
		OutputPath: ".env",
	})
	if cfg.VaultAddr != "http://localhost:8200" {
		t.Errorf("unexpected addr: %s", cfg.VaultAddr)
	}
	if cfg.VaultToken != "s.test" {
		t.Errorf("unexpected token: %s", cfg.VaultToken)
	}
}

func TestFromEnvOverridesWithEnvVar(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://vault.example.com")
	t.Setenv("VAULT_TOKEN", "s.fromenv")

	cfg := preflight.FromEnv(preflight.Config{
		VaultAddr:  "http://localhost:8200",
		VaultToken: "s.old",
	})
	if cfg.VaultAddr != "http://vault.example.com" {
		t.Errorf("expected env override, got %s", cfg.VaultAddr)
	}
	if cfg.VaultToken != "s.fromenv" {
		t.Errorf("expected env override, got %s", cfg.VaultToken)
	}
}

func TestConfigBuildReturnsRunner(t *testing.T) {
	cfg := preflight.Config{
		VaultAddr:  "http://localhost:8200",
		VaultToken: "s.test",
		OutputPath: ".env",
	}
	r := cfg.Build()
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
	// All fields valid — Run should succeed (output dir is '.' which is writable).
	if err := r.Run(); err != nil {
		t.Fatalf("unexpected preflight failure: %v", err)
	}
}
