package pipeline_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpull/internal/pipeline"
)

// stubStage is a test double that applies a simple transformation.
type stubStage struct {
	fn  func(map[string]string) (map[string]string, error)
}

func (s stubStage) Apply(m map[string]string) (map[string]string, error) {
	return s.fn(m)
}

func addKey(k, v string) stubStage {
	return stubStage{fn: func(m map[string]string) (map[string]string, error) {
		m[k] = v
		return m, nil
	}}
}

func failStage() stubStage {
	return stubStage{fn: func(m map[string]string) (map[string]string, error) {
		return nil, errors.New("stage error")
	}}
}

func TestApplyRunsAllStages(t *testing.T) {
	p := pipeline.New(addKey("A", "1"), addKey("B", "2"))
	out, err := p.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestApplyPreservesInputMap(t *testing.T) {
	input := map[string]string{"ORIG": "yes"}
	p := pipeline.New(addKey("NEW", "val"))
	out, _ := p.Apply(input)
	if out["ORIG"] != "yes" {
		t.Errorf("original key missing from output")
	}
	if _, ok := input["NEW"]; ok {
		t.Errorf("Apply mutated the original input map")
	}
}

func TestApplyStopsOnError(t *testing.T) {
	p := pipeline.New(addKey("A", "1"), failStage(), addKey("B", "2"))
	_, err := p.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestApplyEmptyPipelineIsNoop(t *testing.T) {
	p := pipeline.New()
	input := map[string]string{"K": "V"}
	out, err := p.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "V" {
		t.Errorf("expected passthrough, got %v", out)
	}
}
