package scope_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/scope"
)

func TestAllowNoRulesPermitsAll(t *testing.T) {
	s, err := scope.New(scope.Policy{})
	if err != nil {
		t.Fatal(err)
	}
	if !s.Allow("secret/app/db") {
		t.Error("expected path to be allowed")
	}
}

func TestAllowMatchesPrefix(t *testing.T) {
	s, _ := scope.New(scope.Policy{Allow: []string{"secret/app"}})
	if !s.Allow("secret/app/db") {
		t.Error("expected allowed")
	}
	if s.Allow("secret/other/db") {
		t.Error("expected denied")
	}
}

func TestDenyTakesPrecedence(t *testing.T) {
	s, _ := scope.New(scope.Policy{
		Allow: []string{"secret/app"},
		Deny:  []string{"secret/app/internal"},
	})
	if !s.Allow("secret/app/public") {
		t.Error("expected allowed")
	}
	if s.Allow("secret/app/internal/key") {
		t.Error("expected denied")
	}
}

func TestFilterReturnsScopedPaths(t *testing.T) {
	s, _ := scope.New(scope.Policy{Allow: []string{"secret/app"}})
	input := []string{"secret/app/db", "secret/infra/cert", "secret/app/api"}
	got := s.Filter(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(got))
	}
}

func TestValidateRejectsBlankAllow(t *testing.T) {
	_, err := scope.New(scope.Policy{Allow: []string{"", "secret/app"}})
	if err == nil {
		t.Error("expected error for blank allow entry")
	}
}

func TestValidateRejectsConflict(t *testing.T) {
	_, err := scope.New(scope.Policy{
		Allow: []string{"secret/app"},
		Deny:  []string{"secret/app"},
	})
	if err == nil {
		t.Error("expected error for allow/deny conflict")
	}
}
