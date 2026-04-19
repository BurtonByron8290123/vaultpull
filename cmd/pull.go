package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/your-org/vaultpull/internal/vault"
	"github.com/your-org/vaultpull/internal/envwriter"
)

var pullCmd = &cobra.Command{
	Use:   "pull [secret-path]",
	Short: "Pull secrets from Vault and write to .env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().Bool("rotate", false, "Rotate existing secrets before writing")
	pullCmd.Flags().Bool("dry-run", false, "Print secrets to stdout without writing to file")
	pullCmd.Flags().String("prefix", "", "Optional prefix to filter secret keys")
}

func runPull(cmd *cobra.Command, args []string) error {
	secretPath := args[0]

	addr := viper.GetString("vault_addr")
	token := viper.GetString("vault_token")
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("vault token is required: set --vault-token or VAULT_TOKEN")
	}

	client, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	secrets, err := client.ReadSecrets(secretPath)
	if err != nil {
		return fmt.Errorf("failed to read secrets from %q: %w", secretPath, err)
	}

	if len(secrets) == 0 {
		return fmt.Errorf("no secrets found at path %q", secretPath)
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		for k, v := range secrets {
			fmt.Printf("%s=%s\n", k, v)
		}
		return nil
	}

	output := viper.GetString("output")
	prefix, _ := cmd.Flags().GetString("prefix")

	writer := envwriter.New(output)
	if err := writer.Write(secrets, prefix); err != nil {
		return fmt.Errorf("failed to write secrets to %q: %w", output, err)
	}

	fmt.Printf("Successfully wrote %d secrets to %s\n", len(secrets), output)
	return nil
}
