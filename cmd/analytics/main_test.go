package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendEventPostsJSONWithBearerToken(t *testing.T) {
	body := []byte(`{"type":"technical","payload":{"job":"build"}}`)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q, want bearer token", got)
		}

		var event map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			t.Errorf("decode request body: %v", err)
			return
		}
		if event["type"] != "technical" {
			t.Errorf("event type = %v, want technical", event["type"])
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	if err := sendEvent(server.URL, "test-token", body, server.Client()); err != nil {
		t.Fatalf("sendEvent returned an error: %v", err)
	}
}

func TestSendEventReportsHTTPErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = io.WriteString(w, "analytics unavailable")
	}))
	defer server.Close()

	err := sendEvent(server.URL, "", []byte(`{}`), server.Client())
	if err == nil {
		t.Fatal("sendEvent returned nil for an unsuccessful response")
	}
	if !strings.Contains(err.Error(), "502") {
		t.Fatalf("error = %q, want HTTP status", err)
	}
	if strings.Contains(err.Error(), "analytics unavailable") {
		t.Fatalf("error exposes response body: %q", err)
	}
}
