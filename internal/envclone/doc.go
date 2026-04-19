// Package envclone provides a Cloner that copies a map of Vault secrets
// into a destination .env file on disk.
//
// By default existing keys are preserved; set Overwrite=true to replace them.
// DryRun mode reports how many keys would be written without touching the file.
//
// Example:
//
//	p := envclone.DefaultPolicy()
//	p.Overwrite = true
//	c, err := envclone.New(p)
//	n, err := c.Clone(".env", secrets)
package envclone
