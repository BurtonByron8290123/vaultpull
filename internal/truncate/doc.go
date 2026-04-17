// Package truncate provides a Policy type that limits the length of secret
// values before they are written to .env files.
//
// Usage:
//
//	p := truncate.Policy{MaxLen: 256, Suffix: "..."}
//	if err := p.Validate(); err != nil {
//		log.Fatal(err)
//	}
//	safeSecrets := p.ApplyMap(secrets)
//
// A MaxLen of 0 (the default) disables truncation entirely so existing
// behaviour is preserved unless the caller opts in.
package truncate
