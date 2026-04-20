// Package envpatch provides a lightweight patching layer for env maps.
//
// A patch is an ordered list of Ops. Each Op is either a "set" (add or
// overwrite a key/value pair) or a "delete" (remove a key). Ops are applied
// in sequence to an immutable copy of the base map so the original is never
// mutated.
//
// Example:
//
//	p, _ := envpatch.New(envpatch.DefaultPolicy())
//	patched, result, err := p.Apply(base, []envpatch.Op{
//	    {Kind: envpatch.OpSet,    Key: "NEW_KEY",  Value: "hello"},
//	    {Kind: envpatch.OpDelete, Key: "OLD_KEY"},
//	})
package envpatch
