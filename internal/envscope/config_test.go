package envscope

import (
	"testing"
)

func TestFromEnvUsesDefaults(t *testing.T) {
	t.Setenv(envAllow, "")
	t.Setenv(envDeny, "")
	p := FromEnv()
	if len(p.Allow) != 0 || len(p.Deny) != 0 {
		t.Fatalf("expected empty policy, got allow=%v deny=%v", p.Allow, p.Deny)
	}
}

func TestFromEnvReadsAllow(t *testing.T) {
	t.Setenv(envAllow, "APP_, SVC_")
	t.Setenv(envDeny, "")
	p := FromEnv()
	if len(p.Allow) != 2 {
		t.Fatalf("expected 2 allow prefixes, got %v", p.Allow)
	}
	if p.Allow[0] != "APP_" || p.Allow[1] != "SVC_" {
		t.Errorf("unexpected allow values: %v", p.Allow)
	}
}

func TestFromEnvReadsDeny(t *testing.T) {
	t.Setenv(envAllow, "")
	t.Setenv(envDeny, "INTERNAL_,DEBUG_")
	p := FromEnv()
	if len(p.Deny) != 2 {
		t.Fatalf("expected 2 deny prefixes, got %v", p.Deny)
	}
}

func TestFromEnvSkipsBlanks(t *testing.T) {
	t.Setenv(envAllow, "APP_,,SVC_")
	t.Setenv(envDeny, "")
	p := FromEnv()
	if len(p.Allow) != 2 {
		t.Fatalf("expected 2 allow prefixes after blank skip, got %v", p.Allow)
	}
}
