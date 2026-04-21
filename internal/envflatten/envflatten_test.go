package envflatten_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envflatten"
)

func TestFlattenSimpleMap(t *testing.T) {
	f, err := envflatten.New(envflatten.DefaultPolicy())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := map[string]any{
		"db_host": "localhost",
		"db_port": 5432,
	}
	out := f.Flatten(input)
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
}

func TestFlattenNestedMap(t *testing.T) {
	f, _ := envflatten.New(envflatten.DefaultPolicy())
	input := map[string]any{
		"database": map[string]any{
			"host": "db.internal",
			"credentials": map[string]any{
				"user": "admin",
			},
		},
	}
	out := f.Flatten(input)
	if out["DATABASE_HOST"] != "db.internal" {
		t.Errorf("expected DATABASE_HOST=db.internal, got %q", out["DATABASE_HOST"])
	}
	if out["DATABASE_CREDENTIALS_USER"] != "admin" {
		t.Errorf("expected DATABASE_CREDENTIALS_USER=admin, got %q", out["DATABASE_CREDENTIALS_USER"])
	}
}

func TestFlattenWithPrefix(t *testing.T) {
	p := envflatten.DefaultPolicy()
	p.Prefix = "APP"
	f, _ := envflatten.New(p)
	input := map[string]any{"secret": "value"}
	out := f.Flatten(input)
	if out["APP_SECRET"] != "value" {
		t.Errorf("expected APP_SECRET=value, got %q", out["APP_SECRET"])
	}
}

func TestFlattenLowerCaseKeys(t *testing.T) {
	p := envflatten.Policy{Separator: "_", UpperCase: false}
	f, _ := envflatten.New(p)
	input := map[string]any{"MyKey": "val"}
	out := f.Flatten(input)
	if out["MyKey"] != "val" {
		t.Errorf("expected MyKey=val, got %v", out)
	}
}

func TestFlattenCustomSeparator(t *testing.T) {
	p := envflatten.Policy{Separator: ".", UpperCase: false}
	f, _ := envflatten.New(p)
	input := map[string]any{
		"a": map[string]any{"b": "c"},
	}
	out := f.Flatten(input)
	if out["a.b"] != "c" {
		t.Errorf("expected a.b=c, got %v", out)
	}
}

func TestNewRejectsEmptySeparator(t *testing.T) {
	_, err := envflatten.New(envflatten.Policy{Separator: ""})
	if err == nil {
		t.Fatal("expected error for empty separator, got nil")
	}
}

func TestFlattenStringStringMap(t *testing.T) {
	f, _ := envflatten.New(envflatten.DefaultPolicy())
	input := map[string]any{
		"creds": map[string]string{
			"token": "abc123",
		},
	}
	out := f.Flatten(input)
	if out["CREDS_TOKEN"] != "abc123" {
		t.Errorf("expected CREDS_TOKEN=abc123, got %q", out["CREDS_TOKEN"])
	}
}
