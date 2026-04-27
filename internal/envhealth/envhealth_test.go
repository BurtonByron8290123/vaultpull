package envhealth_test

import (
	"testing"

	"github.com/vaultpull/internal/envhealth"
)

func newChecker(p envhealth.Policy) *envhealth.Checker {
	return envhealth.New(p)
}

func TestCheckPassesWhenEnvIsHealthy(t *testing.T) {
	c := newChecker(envhealth.Policy{
		RequiredKeys: []string{"DB_URL", "API_KEY"},
	})
	report := c.Check(map[string]string{"DB_URL": "postgres://localhost", "API_KEY": "secret"})
	if report.Status != envhealth.StatusOK {
		t.Fatalf("expected OK, got %v violations: %v", report.Status, report.Violations)
	}
}

func TestCheckFailsWhenRequiredKeyMissing(t *testing.T) {
	c := newChecker(envhealth.Policy{RequiredKeys: []string{"MISSING_KEY"}})
	report := c.Check(map[string]string{})
	if report.Status != envhealth.StatusError {
		t.Fatal("expected StatusError")
	}
	if len(report.Violations) != 1 || report.Violations[0].Key != "MISSING_KEY" {
		t.Fatalf("unexpected violations: %v", report.Violations)
	}
}

func TestCheckFailsWhenRequiredKeyEmpty(t *testing.T) {
	c := newChecker(envhealth.Policy{RequiredKeys: []string{"TOKEN"}})
	report := c.Check(map[string]string{"TOKEN": "   "})
	if report.Status != envhealth.StatusError {
		t.Fatal("expected StatusError for empty value")
	}
}

func TestCheckFailsWhenForbiddenKeyPresent(t *testing.T) {
	c := newChecker(envhealth.Policy{ForbiddenKeys: []string{"DEBUG"}})
	report := c.Check(map[string]string{"DEBUG": "true"})
	if report.Status != envhealth.StatusError {
		t.Fatal("expected StatusError for forbidden key")
	}
	if report.Violations[0].Key != "DEBUG" {
		t.Fatalf("expected violation for DEBUG, got %v", report.Violations)
	}
}

func TestCheckNoEmptyValuesFlag(t *testing.T) {
	c := newChecker(envhealth.Policy{NoEmptyValues: true})
	report := c.Check(map[string]string{"KEY_A": "value", "KEY_B": ""})
	if report.Status != envhealth.StatusError {
		t.Fatal("expected StatusError when empty value present")
	}
	if len(report.Violations) != 1 || report.Violations[0].Key != "KEY_B" {
		t.Fatalf("unexpected violations: %v", report.Violations)
	}
}

func TestCheckNoEmptyValuesFlagPassesWhenAllFilled(t *testing.T) {
	c := newChecker(envhealth.Policy{NoEmptyValues: true})
	report := c.Check(map[string]string{"A": "1", "B": "2"})
	if report.Status != envhealth.StatusOK {
		t.Fatalf("expected OK, got %v", report.Violations)
	}
}

func TestSummaryOK(t *testing.T) {
	r := envhealth.Report{Status: envhealth.StatusOK}
	if r.Summary() != "env healthy: no issues found" {
		t.Fatalf("unexpected summary: %s", r.Summary())
	}
}

func TestSummaryError(t *testing.T) {
	r := envhealth.Report{
		Status:     envhealth.StatusError,
		Violations: []envhealth.Violation{{Key: "X", Message: "missing"}},
	}
	s := r.Summary()
	if s == "" {
		t.Fatal("summary should not be empty")
	}
}

func TestDefaultPolicyProducesNoViolations(t *testing.T) {
	c := envhealth.New(envhealth.DefaultPolicy())
	report := c.Check(map[string]string{"ANYTHING": ""})
	if report.Status != envhealth.StatusOK {
		t.Fatalf("default policy should pass everything, got %v", report.Violations)
	}
}
