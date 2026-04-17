// Package redact provides key-based value redaction for vaultpull.
//
// It is used to prevent sensitive secret values from appearing in
// terminal output, audit logs, and diff views. Redaction is based on
// a configurable list of key name patterns (case-insensitive).
//
// Usage:
//
//	r := redact.FromConfig(cfg.Redact)
//	safeMap := r.Map(secrets)
//
// The default sensitive key list covers common patterns such as
// PASSWORD, TOKEN, SECRET, API_KEY and more. Additional keys can be
// supplied via configuration or the VAULTPULL_REDACT_KEYS environment
// variable (comma-separated).
package redact
