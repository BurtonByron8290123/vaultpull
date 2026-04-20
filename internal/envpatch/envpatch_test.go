package envpatch_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envpatch"
)

func newPatcher(t *testing.T) *envpatch.Patcher {
	t.Helper()
	p, err := envpatch.New(envpatch.DefaultPolicy())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return p
}

func TestApplyAddsNewKey(t *testing.T) {
	patcher := newPatcher(t)
	base := map[string]string{"A": "1"}
	out, res, err := patcher.Apply(base, []envpatch.Op{
		{Kind: envpatch.OpSet, Key: "B", Value: "2"},
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %q", out["B"])
	}
	if res.Added != 1 || res.Updated != 0 {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestApplyUpdatesExistingKey(t *testing.T) {
	patcher := newPatcher(t)
	base := map[string]string{"A": "old"}
	out, res, err := patcher.Apply(base, []envpatch.Op{
		{Kind: envpatch.OpSet, Key: "A", Value: "new"},
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if out["A"] != "new" {
		t.Errorf("expected A=new, got %q", out["A"])
	}
	if res.Updated != 1 {
		t.Errorf("expected 1 update, got %d", res.Updated)
	}
}

func TestApplyDeletesKey(t *testing.T) {
	patcher := newPatcher(t)
	base := map[string]string{"A": "1", "B": "2"}
	out, res, err := patcher.Apply(base, []envpatch.Op{
		{Kind: envpatch.OpDelete, Key: "A"},
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if _, ok := out["A"]; ok {
		t.Error("expected A to be deleted")
	}
	if res.Deleted != 1 {
		t.Errorf("expected 1 deletion, got %d", res.Deleted)
	}
}

func TestApplyBaseNotMutated(t *testing.T) {
	patcher := newPatcher(t)
	base := map[string]string{"A": "1"}
	_, _, err := patcher.Apply(base, []envpatch.Op{
		{Kind: envpatch.OpSet, Key: "A", Value: "999"},
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if base["A"] != "1" {
		t.Errorf("base was mutated: A=%q", base["A"])
	}
}

func TestApplyIgnoreExistingSkipsUpdate(t *testing.T) {
	p, _ := envpatch.New(envpatch.Policy{IgnoreExisting: true})
	base := map[string]string{"A": "original"}
	out, res, err := p.Apply(base, []envpatch.Op{
		{Kind: envpatch.OpSet, Key: "A", Value: "overwrite"},
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if out["A"] != "original" {
		t.Errorf("expected A=original, got %q", out["A"])
	}
	if res.Updated != 0 || res.Added != 0 {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestApplyUnknownOpReturnsError(t *testing.T) {
	patcher := newPatcher(t)
	_, _, err := patcher.Apply(nil, []envpatch.Op{
		{Kind: "invalid", Key: "X"},
	})
	if err == nil {
		t.Error("expected error for unknown op kind")
	}
}

func TestApplyEmptyKeyReturnsError(t *testing.T) {
	patcher := newPatcher(t)
	_, _, err := patcher.Apply(nil, []envpatch.Op{
		{Kind: envpatch.OpSet, Key: "", Value: "v"},
	})
	if err == nil {
		t.Error("expected error for empty key")
	}
}
