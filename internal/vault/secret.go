package vault

import (
	"context"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Secret represents a set of key-value pairs fetched from a Vault path.
type Secret struct {
	Path string
	Data map[string]string
}

// FetchSecrets reads key-value pairs from the given Vault KV path.
// It supports both KV v1 and KV v2 by detecting the "data" wrapper.
func FetchSecrets(client *vaultapi.Client, path string) (*Secret, error) {
	return fetchSecrets(client, path)
}

func fetchSecrets(client *vaultapi.Client, path string) (*Secret, error) {
	ctx := context.Background()

	secret, err := client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("vault read %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("vault path %q returned no data", path)
	}

	raw := secret.Data

	// KV v2 wraps actual data under a "data" key.
	if nested, ok := raw["data"]; ok {
		if nestedMap, ok := nested.(map[string]interface{}); ok {
			raw = nestedMap
		}
	}

	result := &Secret{
		Path: path,
		Data: make(map[string]string, len(raw)),
	}

	for k, v := range raw {
		key := strings.ToUpper(k)
		switch val := v.(type) {
		case string:
			result.Data[key] = val
		case fmt.Stringer:
			result.Data[key] = val.String()
		default:
			result.Data[key] = fmt.Sprintf("%v", val)
		}
	}

	return result, nil
}
