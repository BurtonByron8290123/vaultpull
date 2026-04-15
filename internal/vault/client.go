package vault

import (
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	api *vaultapi.Client
}

// NewClient creates a new Vault client with the given address and token.
func NewClient(addr, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("initializing vault api client: %w", err)
	}

	api.SetToken(token)
	return &Client{api: api}, nil
}

// ReadSecrets reads key-value secrets from the given Vault path.
// Supports both KV v1 and KV v2 (data/ prefix) paths.
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading path %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	data, ok := secret.Data["data"]
	if !ok {
		// KV v1 — data is at top level
		data = secret.Data
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected secret data format at path %q", path)
	}

	result := make(map[string]string, len(dataMap))
	for k, v := range dataMap {
		strVal, ok := v.(string)
		if !ok {
			strVal = fmt.Sprintf("%v", v)
		}
		result[k] = strVal
	}

	return result, nil
}
