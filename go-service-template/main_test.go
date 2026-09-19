package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestServer wires the same routes as main() so the handlers can be exercised
// without binding a port. It keeps the service template's coverage above the
// platform's minimum (conformance rule test-coverage, ADR atlas-api/0019).
func newTestServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/dbz", dbzHandler)
	mux.HandleFunc("/", rootHandler)
	return mux
}

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("body = %q, want %q", rec.Body.String(), "ok")
	}
}

func TestRootWithoutDatabase(t *testing.T) {
	t.Setenv("DB_HOST", "")
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "database:") {
		t.Fatalf("body should not mention a database when DB_HOST is empty: %q", rec.Body.String())
	}
}

func TestRootWithDatabase(t *testing.T) {
	t.Setenv("DB_HOST", "postgres.example")
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(rec.Body.String(), "postgres.example") {
		t.Fatalf("body should mention the database host: %q", rec.Body.String())
	}
}

func TestDbzNoDatabase(t *testing.T) {
	t.Setenv("DB_HOST", "")
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/dbz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "no database configured") {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestDbzUnreachable(t *testing.T) {
	// Reserved TEST-NET-1 address (RFC 5737): the dial fails fast without touching the network.
	t.Setenv("DB_HOST", "192.0.2.1")
	t.Setenv("DB_PORT", "5432")
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/dbz", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "unreachable") {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestDbzReachable(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	host, port, _ := net.SplitHostPort(ln.Addr().String())
	t.Setenv("DB_HOST", host)
	t.Setenv("DB_PORT", port)
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/dbz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "reachable") {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestGetenv(t *testing.T) {
	t.Setenv("SAMPLE", "")
	if got := getenv("SAMPLE", "fallback"); got != "fallback" {
		t.Fatalf("getenv empty = %q, want fallback", got)
	}
	t.Setenv("SAMPLE", "value")
	if got := getenv("SAMPLE", "fallback"); got != "value" {
		t.Fatalf("getenv set = %q, want value", got)
	}
}

func TestDial(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	if err := dial(ln.Addr().String()); err != nil {
		t.Fatalf("dial reachable: %v", err)
	}
	if err := dial("192.0.2.1:5432"); err == nil {
		t.Fatal("dial to unreachable address should fail")
	}
}
