package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpull is a CLI tool that pulls secrets from HashiCorp Vault
and writes them into local .env files with optional rotation support.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .vaultpull.yaml)")
	rootCmd.PersistentFlags().String("vault-addr", "http://127.0.0.1:8200", "Vault server address")
	rootCmd.PersistentFlags().String("vault-token", "", "Vault token (overrides VAULT_TOKEN env var)")
	rootCmd.PersistentFlags().String("output", ".env", "Output .env file path")

	_ = viper.BindPFlag("vault_addr", rootCmd.PersistentFlags().Lookup("vault-addr"))
	_ = viper.BindPFlag("vault_token", rootCmd.PersistentFlags().Lookup("vault-token"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".vaultpull")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("VAULTPULL")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
