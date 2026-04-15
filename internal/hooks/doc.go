// Package hooks implements pre-pull and post-pull lifecycle hook execution
// for vaultpull. Hooks are arbitrary shell commands that run before or after
// secrets are fetched from Vault and written to the local .env file.
//
// Example configuration (vaultpull.yaml):
//
//	hooks:
//	  pre_pull: "./scripts/backup-env.sh"
//	  post_pull: "systemctl reload myapp"
//	  timeout: 60s
//
// A non-zero exit code from a hook causes vaultpull to abort with an error.
// If a hook exceeds the configured timeout it is killed and an error is
// returned. The default timeout is 30 seconds.
package hooks
