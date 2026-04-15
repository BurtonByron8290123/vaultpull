// Package prompt provides utilities for interactive terminal prompts used
// by vaultpull commands.
//
// Currently it exposes a Confirmer that renders yes/no questions and parses
// the user's response.  Commands that perform destructive or irreversible
// operations (e.g. overwriting .env files, pruning backups) should gate
// those actions behind a confirmation when running in interactive mode.
//
// Example:
//
//	c := prompt.New()
//	ok, err := c.Ask("Overwrite existing .env file?", false)
//	if err != nil || !ok {
//		return
//	}
//	// proceed with write
package prompt
