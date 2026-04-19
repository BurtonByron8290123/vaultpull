// Package scope provides path-level access control for Vault secret pulls.
//
// A Policy declares which path prefixes are allowed or denied. The Scope
// enforcer applies the policy so that only permitted paths are fetched,
// preventing accidental exposure of secrets outside the intended namespace.
//
// Deny rules always take precedence over allow rules. When no allow prefixes
// are configured every path is implicitly allowed (subject to deny rules).
//
// Policies can be constructed programmatically or loaded from the environment
// via FromEnv using VAULTPULL_SCOPE_ALLOW and VAULTPULL_SCOPE_DENY.
package scope
