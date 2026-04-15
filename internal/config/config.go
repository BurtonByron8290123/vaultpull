package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds all runtime configuration for vaultpull.
type Config struct {
	// Vault connection
	VaultAddr  string `mapstructure:"vault_addr"`
	VaultToken string `mapstructure:"vault_token"`
	VaultMount string `mapstructure:"vault_mount"`

	// Secret paths to pull (may contain template variables)
	Paths []string `mapstructure:"paths"`

	// Output
	OutputFile string `mapstructure:"output_file"`
	Merge      bool   `mapstructure:"merge"`

	// Rotation
	BackupDir  string `mapstructure:"backup_dir"`
	MaxBackups int    `mapstructure:"max_backups"`

	// Audit
	AuditLog string `mapstructure:"audit_log"`

	// Template variables used to expand Paths
	TemplateVars map[string]string `mapstructure:"template_vars"`

	// Timeout for Vault HTTP requests
	Timeout time.Duration `mapstructure:"timeout"`
}

// Load reads configuration from file and environment, then validates it.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	v.SetDefault("vault_addr", "http://127.0.0.1:8200")
	v.SetDefault("vault_mount", "secret")
	v.SetDefault("output_file", ".env")
	v.SetDefault("merge", true)
	v.SetDefault("backup_dir", ".env.backups")
	v.SetDefault("max_backups", 5)
	v.SetDefault("timeout", 10*time.Second)

	v.SetEnvPrefix("VAULTPULL")
	v.AutomaticEnv()

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName(".vaultpull")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath(os.ExpandEnv("$HOME"))
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && cfgFile != "" {
			return nil, fmt.Errorf("config: read %q: %w", cfgFile, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validate(c *Config) error {
	if c.VaultToken == "" {
		return errors.New("config: vault_token is required (set VAULTPULL_VAULT_TOKEN or vault_token in config)")
	}
	if len(c.Paths) == 0 {
		return errors.New("config: at least one path is required")
	}
	return nil
}
