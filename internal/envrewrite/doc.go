// Package envrewrite provides regex-based value rewriting for env maps.
//
// Rules are applied in declaration order. Each rule may optionally be
// scoped to keys that share a common prefix via the key_glob field.
//
// Usage:
//
//	p := envrewrite.Policy{
//	    Rules: []envrewrite.Rule{
//	        {Pattern: `^http://`, Replacement: "https://"},
//	    },
//	}
//	rw, err := envrewrite.New(p)
//	if err != nil { ... }
//	result := rw.Apply(secretMap)
package envrewrite
