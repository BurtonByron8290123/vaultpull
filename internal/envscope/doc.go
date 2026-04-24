// Package envscope provides prefix-based allow/deny filtering for env var maps.
//
// It is intended to be used as a pipeline stage that restricts which secrets
// fetched from Vault are written to a given .env file. Rules follow the same
// semantics as internal/scope: deny takes precedence over allow, and an empty
// allow list permits all keys not explicitly denied.
//
// Configuration can be loaded from environment variables via FromEnv or
// supplied programmatically via a Policy struct.
package envscope
