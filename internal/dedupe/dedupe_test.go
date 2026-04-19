package dedupe_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/dedupe"
)

func TestDefaultPolicyLastWins(t *testing.T) {
	d, err := dedupe.New(dedupe.DefaultPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	out, err := d.Merge(a, b)
	if err != nil {
		t.Fatalf("merge error: %v", err)
	}
	if out["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", out["KEY"])
	}
}

func TestMergeNoConflict(t *testing.T) {
	d, _ := dedupe.New(dedupe.DefaultPolicy())
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	out, err := d.Merge(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Errorf("unexpected result: %v", out)
	}
}

func TestFailOnConflict(t *testing.T) {
	d, err := dedupe.New(dedupe.Policy{FailOnConflict: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a := map[string]string{"KEY": "v1"}
	b := map[string]string{"KEY": "v2"}
	_, err = d.Merge(a, b)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}

func TestFailOnConflictSameValueNoError(t *testing.T) {
	d, _ := dedupe.New(dedupe.Policy{FailOnConflict: true})
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same"}
	_, err := d.Merge(a, b)
	if err != nil {
		t.Fatalf("expected no error for identical values, got: %v", err)
	}
}

func TestMutuallyExclusivePolicyReturnsError(t *testing.T) {
	_, err := dedupe.New(dedupe.Policy{LastWins: true, FailOnConflict: true})
	if err == nil {
		t.Fatal("expected error for conflicting policy options")
	}
}

func TestMergeEmptyMaps(t *testing.T) {
	d, _ := dedupe.New(dedupe.DefaultPolicy())
	out, err := d.Merge(map[string]string{}, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
