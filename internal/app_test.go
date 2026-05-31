package internal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3" // driver for test DB connection in New()
	"scavo/internal"
)

func TestServer_WithSecret_BuildsAndServesHealth(t *testing.T) {
	t.Setenv("SESSION_SECRET", "ci-test-secret-very-long-and-random-123456")

	handler, addr := internal.New()
	if handler == nil {
		t.Fatal("New() returned nil handler even with SESSION_SECRET set")
	}
	if addr == "" {
		t.Error("expected non-empty listen address")
	}

	// Exercise a real handler
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET /health = %d, want 200", rec.Code)
	}
	if rec.Body.String() != "OK" {
		t.Errorf("health body = %q, want OK", rec.Body.String())
	}
}
