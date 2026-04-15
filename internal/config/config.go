package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	VaultAddr  string            `mapstructure:"vault_addr"`
	VaultToken string            `mapstructure:"vault_token"`
	VaultPath  string            `mapstructure:"vault_path"`
	OutputFile string            `mapstructure:"output_file"`
	Rotate     bool              `mapstructure:"rotate"`
	BackupDir  string            `mapstructure:"backup_dir"`
	Mappings   map[string]string `mapstructure:"mappings"`
}

// Load reads configuration from file and environment variables
func Load(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".vaultpull")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("VAULTPULL")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("output_file", ".env")
	viper.SetDefault("backup_dir", ".env.backups")
	viper.SetDefault("rotate", false)
	viper.SetDefault("vault_addr", "http://127.0.0.1:8200")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Allow VAULT_TOKEN from environment directly
	if token := os.Getenv("VAULT_TOKEN"); token != "" && viper.GetString("vault_token") == "" {
		viper.Set("vault_token", token)
	}

	if addr := os.Getenv("VAULT_ADDR"); addr != "" && viper.GetString("vault_addr") == "http://127.0.0.1:8200" {
		viper.Set("vault_addr", addr)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.VaultAddr == "" {
		return fmt.Errorf("vault_addr is required")
	}
	if c.VaultToken == "" {
		return fmt.Errorf("vault_token is required (set VAULT_TOKEN or vault_token in config)")
	}
	if c.VaultPath == "" {
		return fmt.Errorf("vault_path is required")
	}
	return nil
}
