// Package pipeline chains secret post-processing steps (sanitize, transform,
// redact, truncate) into a single Apply call used by the pull command.
package pipeline

import (
	"fmt"

	"github.com/your-org/vaultpull/internal/redact"
	"github.com/your-org/vaultpull/internal/sanitize"
	"github.com/your-org/vaultpull/internal/transform"
	"github.com/your-org/vaultpull/internal/truncate"
)

// Stage is a single processing step.
type Stage interface {
	Apply(map[string]string) (map[string]string, error)
}

// Pipeline runs secrets through an ordered list of stages.
type Pipeline struct {
	stages []Stage
}

// New builds a Pipeline from the supplied stages.
func New(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Default constructs a Pipeline with the standard set of stages using
// package-level default policies.
func Default() *Pipeline {
	san := sanitize.New(sanitize.DefaultPolicy())
	tr := transform.New(transform.Policy{})
	rd := redact.New(redact.DefaultSensitiveKeys, "***")
	trc := truncate.New(truncate.DefaultPolicy())
	return New(san, tr, rd, trc)
}

// Apply runs secrets through every stage in order, returning the final map.
func (p *Pipeline) Apply(secrets map[string]string) (map[string]string, error) {
	current := copyMap(secrets)
	for i, s := range p.stages {
		var err error
		current, err = s.Apply(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline stage %d: %w", i, err)
		}
	}
	return current, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
