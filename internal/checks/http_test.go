package checks

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckHTTPSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := CheckHTTP(server.URL, time.Second)
	if !result.Success {
		t.Fatalf("expected success, got failure: %s", result.Message)
	}
}

func TestCheckHTTPFailureStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	result := CheckHTTP(server.URL, time.Second)
	if result.Success {
		t.Fatalf("expected failure for 503 status")
	}
}
