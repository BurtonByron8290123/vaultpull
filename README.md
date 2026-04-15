# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files with rotation support.

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Authenticate with your Vault instance and run `vaultpull` to sync secrets into a local `.env` file.

```bash
# Set your Vault address and token
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

# Pull secrets from a Vault path into a .env file
vaultpull pull --path secret/data/myapp --output .env

# Pull and rotate secrets (re-generates dynamic credentials)
vaultpull pull --path secret/data/myapp --output .env --rotate

# Specify a custom output file and Vault namespace
vaultpull pull --path secret/data/myapp --output config/.env.local --namespace my-team
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to pull from | *(required)* |
| `--output` | Output `.env` file path | `.env` |
| `--rotate` | Rotate secrets before syncing | `false` |
| `--namespace` | Vault namespace | `""` |

---

## Configuration

`vaultpull` respects standard Vault environment variables:

- `VAULT_ADDR` — Vault server address
- `VAULT_TOKEN` — Authentication token
- `VAULT_NAMESPACE` — Vault namespace (optional)

---

## License

[MIT](LICENSE) © 2024 yourusername