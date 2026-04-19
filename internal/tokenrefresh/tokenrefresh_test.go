package tokenrefresh

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestServer(ttlSec float64, renewOK bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/token/lookup-self":
			body := map[string]interface{}{"data": map[string]interface{}{"ttl": ttlSec}}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
		case "/v1/auth/token/renew-self":
			if renewOK {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{"auth": map[string]interface{}{}})
			} else {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]interface{}{"errors": []string{"permission denied"}})
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func newClient(t *testing.T, addr string) *vaultapi.Client {
	t.Helper()
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr
	c, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	c.SetToken("test-token")
	return c
}

func TestEnsureValidSkipsRenewalWhenFresh(t *testing.T) {
	srv := newTestServer(600, false)
	defer srv.Close()
	r, _ := New(newClient(t, srv.URL), Policy{RenewThreshold: 5 * time.Minute, MaxRetries: 1})
	if err := r.EnsureValid(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnsureValidRenewsWhenNearExpiry(t *testing.T) {
	srv := newTestServer(60, true)
	defer srv.Close()
	r, _ := New(newClient(t, srv.URL), Policy{RenewThreshold: 5 * time.Minute, MaxRetries: 1})
	if err := r.EnsureValid(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnsureValidReturnsErrorAfterRetries(t *testing.T) {
	srv := newTestServer(60, false)
	defer srv.Close()
	r, _ := New(newClient(t, srv.URL), Policy{RenewThreshold: 5 * time.Minute, MaxRetries: 2})
	if err := r.EnsureValid(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewRejectsInvalidPolicy(t *testing.T) {
	_, err := New(nil, Policy{RenewThreshold: 0, MaxRetries: 1})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewRejectsZeroMaxRetries(t *testing.T) {
	_, err := New(nil, Policy{RenewThreshold: 5 * time.Minute, MaxRetries: 0})
	if err == nil {
		t.Fatal("expected validation error for zero MaxRetries")
	}
}
