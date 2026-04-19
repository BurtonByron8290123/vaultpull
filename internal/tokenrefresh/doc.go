// Package tokenrefresh manages automatic renewal of HashiCorp Vault tokens.
//
// It inspects the remaining TTL of the active token via the lookup-self
// endpoint and triggers a renew-self call whenever the TTL falls below a
// configurable threshold.  Renewal is retried up to MaxRetries times before
// an error is surfaced to the caller.
//
// Configuration can be supplied programmatically via Policy or loaded from
// environment variables with FromEnv:
//
//	VAULTPULL_TOKEN_RENEW_THRESHOLD_SEC  – seconds of remaining TTL that
//	                                        trigger renewal (default 300)
//	VAULTPULL_TOKEN_RENEW_MAX_RETRIES    – maximum renewal attempts (default 3)
package tokenrefresh
