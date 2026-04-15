package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestClient(t *testing.T, handler http.Handler) *vaultapi.Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	cfg := vaultapi.DefaultConfig()
	cfg.Address = ts.URL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	client.SetToken("test-token")
	return client
}

func TestFetchSecretsKVv1(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{"DB_PASS": "secret", "API_KEY": "abc123"},
	}
	client := newTestClient(t, vaultHandler(t, payload))

	s, err := FetchSecrets(client, "secret/myapp")
	if err != nil {
		t. error: %v", err)
	}
	if s.Data["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %q_PASS"])
	}
	if s.Data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", s.Data["API_KEY"])
	}
}

func TestFetchSecretsKeyNormalization(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{"db_pass": "lower"},
	}
	client := newTestClient(t, vaultHandler(t, payload))

	s, err := FetchSecrets(client, "secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.Data["DB_PASS"]; !ok {
		t.Error("expected key DB_PASS to be uppercased")
	}
}

func TestFetchSecretsEmptyPath(t *testing.T) {
	client := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": nil})
	}))

	_, err := FetchSecrets(client, "secret/empty")
	if err == nil {
		t.Error("expected error for nil data, got nil")
	}
}

func vaultHandler(t *testing.T, payload map[string]interface{}) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	})
}
