package pagesize_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/pagesize"
)

func TestDefaultPolicyIsValid(t *testing.T) {
	p := pagesize.DefaultPolicy()
	if err := p.Validate(); err != nil {
		t.Fatalf("expected valid default policy, got: %v", err)
	}
}

func TestValidateRejectsBelowMin(t *testing.T) {
	p := pagesize.Policy{PageSize: 0}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for page size 0")
	}
}

func TestValidateRejectsAboveMax(t *testing.T) {
	p := pagesize.Policy{PageSize: 1001}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for page size 1001")
	}
}

func TestPagesCalculation(t *testing.T) {
	p := pagesize.Policy{PageSize: 10}
	cases := []struct{ total, want int }{
		{0, 0},
		{1, 1},
		{10, 1},
		{11, 2},
		{25, 3},
	}
	for _, c := range cases {
		got := p.Pages(c.total)
		if got != c.want {
			t.Errorf("Pages(%d) = %d, want %d", c.total, got, c.want)
		}
	}
}

func TestSliceReturnsCorrectWindow(t *testing.T) {
	p := pagesize.Policy{PageSize: 3}
	items := []string{"a", "b", "c", "d", "e"}

	got := p.Slice(items, 0)
	if len(got) != 3 || got[0] != "a" {
		t.Fatalf("page 0 unexpected: %v", got)
	}

	got = p.Slice(items, 1)
	if len(got) != 2 || got[0] != "d" {
		t.Fatalf("page 1 unexpected: %v", got)
	}
}

func TestSliceOutOfBoundsReturnsNil(t *testing.T) {
	p := pagesize.Policy{PageSize: 10}
	if got := p.Slice([]string{"a"}, 5); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestFromEnvUsesDefault(t *testing.T) {
	t.Setenv("VAULTPULL_PAGE_SIZE", "")
	p := pagesize.FromEnv()
	if p.PageSize != pagesize.DefaultPageSize {
		t.Fatalf("expected %d, got %d", pagesize.DefaultPageSize, p.PageSize)
	}
}

func TestFromEnvReadsPageSize(t *testing.T) {
	t.Setenv("VAULTPULL_PAGE_SIZE", "50")
	p := pagesize.FromEnv()
	if p.PageSize != 50 {
		t.Fatalf("expected 50, got %d", p.PageSize)
	}
}
