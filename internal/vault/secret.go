package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/your-org/vaultpull/internal/retry"
)

// FetchSecrets retrieves key/value pairs from Vault at the given path.
// It transparently handles both KV v1 and KV v2 mounts and retries
// transient HTTP errors according to the supplied policy.
func FetchSecrets(ctx context.Context, c *Client, path string, p retry.Policy) (map[string]string, error) {
	if path == "" {
		return nil, fmt.Errorf("vault: path must not be empty")
	}
	var result map[string]string
	err := retry.Do(ctx, p, isTransientVaultError, func() error {
		var ferr error
		result, ferr = fetchSecrets(ctx, c, path)
		return ferr
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func fetchSecrets(ctx context.Context, c *Client, path string) (map[string]string, error) {
	secret, err := c.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("vault: read %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault: no secret found at %q", path)
	}

	// KV v2 wraps data under secret.Data["data"].
	data := secret.Data
	if nested, ok := secret.Data["data"]; ok {
		if m, ok := nested.(map[string]interface{}); ok {
			data = m
		}
	}

	out := make(map[string]string, len(data))
	for k, v := range data {
		normKey := strings.ToUpper(strings.ReplaceAll(k, "-", "_"))
		out[normKey] = fmt.Sprintf("%v", v)
	}
	return out, nil
}

// isTransientVaultError returns true for errors that are worth retrying
// (network timeouts, 5xx responses).
func isTransientVaultError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	transientPhrases := []string{
		"connection refused",
		"timeout",
		"EOF",
		"503",
		"502",
		"500",
	}
	for _, phrase := range transientPhrases {
		if strings.Contains(msg, phrase) {
			return true
		}
	}
	return false
}
