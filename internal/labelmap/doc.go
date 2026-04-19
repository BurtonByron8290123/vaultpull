// Package labelmap provides key-renaming for secrets fetched from Vault.
//
// Rules are defined as a YAML list of {from, to} pairs and applied after
// secrets are retrieved but before they are written to the .env file.
//
// Example YAML:
//
//	- from: DB_PASSWORD
//	  to: DATABASE_PASSWORD
//	- from: API_TOKEN
//	  to: SERVICE_API_KEY
//
// Usage:
//
//	m, err := labelmap.LoadConfig("labels.yaml")
//	if err != nil { ... }
//	renamed := m.Apply(secrets)
package labelmap
