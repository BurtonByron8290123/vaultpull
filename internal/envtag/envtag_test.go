package envtag_test

import (
	"testing"

	"github.com/user/vaultpull/internal/envtag"
)

func newTagger(t *testing.T, rules []envtag.Rule) *envtag.Tagger {
	t.Helper()
	tag, err := envtag.New(envtag.Policy{Rules: rules})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tag
}

func TestApplyMatchingRule(t *testing.T) {
	tag := newTagger(t, []envtag.Rule{
		{Prefix: "DB_", Tags: []envtag.Tag{{Key: "group", Value: "database"}}},
	})
	env := map[string]string{"DB_HOST": "localhost", "APP_NAME": "myapp"}
	out := tag.Apply(env)

	if out["DB_HOST"] != "localhost" {
		t.Errorf("original key missing")
	}
	if out["__TAG_GROUP__DB_HOST"] != "database" {
		t.Errorf("tag key missing or wrong value: %q", out["__TAG_GROUP__DB_HOST"])
	}
	if _, ok := out["__TAG_GROUP__APP_NAME"]; ok {
		t.Errorf("tag should not be applied to non-matching key")
	}
}

func TestApplyNoMatchReturnsOriginal(t *testing.T) {
	tag := newTagger(t, []envtag.Rule{
		{Prefix: "SECRET_", Tags: []envtag.Tag{{Key: "sensitivity", Value: "high"}}},
	})
	env := map[string]string{"APP_KEY": "val"}
	out := tag.Apply(env)
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	tag := newTagger(t, []envtag.Rule{
		{Prefix: "X_", Tags: []envtag.Tag{{Key: "env", Value: "prod"}}},
	})
	env := map[string]string{"X_FOO": "bar"}
	_ = tag.Apply(env)
	if len(env) != 1 {
		t.Errorf("input map was mutated")
	}
}

func TestApplyMultipleTags(t *testing.T) {
	tag := newTagger(t, []envtag.Rule{
		{Prefix: "SVC_", Tags: []envtag.Tag{
			{Key: "tier", Value: "backend"},
			{Key: "owner", Value: "platform"},
		}},
	})
	out := tag.Apply(map[string]string{"SVC_URL": "http://example.com"})
	if out["__TAG_TIER__SVC_URL"] != "backend" {
		t.Errorf("tier tag missing")
	}
	if out["__TAG_OWNER__SVC_URL"] != "platform" {
		t.Errorf("owner tag missing")
	}
}

func TestNewRejectsEmptyPrefix(t *testing.T) {
	_, err := envtag.New(envtag.Policy{Rules: []envtag.Rule{
		{Prefix: "", Tags: []envtag.Tag{{Key: "k", Value: "v"}}},
	}})
	if err == nil {
		t.Error("expected error for empty prefix")
	}
}

func TestNewRejectsEmptyTagKey(t *testing.T) {
	_, err := envtag.New(envtag.Policy{Rules: []envtag.Rule{
		{Prefix: "A_", Tags: []envtag.Tag{{Key: "", Value: "v"}}},
	}})
	if err == nil {
		t.Error("expected error for empty tag key")
	}
}

func TestApplyMapSubsetsKeys(t *testing.T) {
	tag := newTagger(t, []envtag.Rule{
		{Prefix: "DB_", Tags: []envtag.Tag{{Key: "group", Value: "db"}}},
	})
	env := map[string]string{"DB_HOST": "h", "DB_PORT": "5432", "APP": "x"}
	out := tag.ApplyMap(env, []string{"DB_HOST"})
	if _, ok := out["DB_PORT"]; ok {
		t.Errorf("DB_PORT should not be in output")
	}
	if out["__TAG_GROUP__DB_HOST"] != "db" {
		t.Errorf("tag not applied to subset key")
	}
}
