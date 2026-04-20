// Package envpatch applies a partial update (patch) to an existing env map,
// supporting add, update, and delete operations in a single pass.
package envpatch

import "fmt"

// OpKind describes the type of patch operation.
type OpKind string

const (
	OpSet    OpKind = "set"    // add or overwrite a key
	OpDelete OpKind = "delete" // remove a key
)

// Op is a single patch operation.
type Op struct {
	Kind  OpKind
	Key   string
	Value string // ignored for OpDelete
}

// Result summarises what changed after Apply.
type Result struct {
	Added   int
	Updated int
	Deleted int
}

// Patcher applies a sequence of Ops to an env map.
type Patcher struct {
	policy Policy
}

// New returns a Patcher using the given policy.
func New(p Policy) (*Patcher, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &Patcher{policy: p}, nil
}

// Apply executes ops against base (mutating a copy) and returns the patched map
// together with a Result summary. base is never modified.
func (p *Patcher) Apply(base map[string]string, ops []Op) (map[string]string, Result, error) {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	var res Result
	for _, op := range ops {
		if op.Key == "" {
			return nil, Result{}, fmt.Errorf("envpatch: op key must not be empty")
		}
		switch op.Kind {
		case OpSet:
			if _, exists := out[op.Key]; exists {
				if p.policy.IgnoreExisting {
					continue
				}
				res.Updated++
			} else {
				res.Added++
			}
			out[op.Key] = op.Value
		case OpDelete:
			if _, exists := out[op.Key]; exists {
				delete(out, op.Key)
				res.Deleted++
			}
		default:
			return nil, Result{}, fmt.Errorf("envpatch: unknown op kind %q", op.Kind)
		}
	}
	return out, res, nil
}
