package hc

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetSharedClient(t *testing.T) {
	t.Helper()

	prev := client
	client = nil

	t.Cleanup(func() {
		client = prev
	})
}

func newTestSettings() *HTTPClientSettings {
	settings := NewClientSettings(true, true, 5, 5)
	settings.ReUseClient = false
	settings.ReturnCert = false
	return settings
}

func TestMakeCallReturnsHeadersAndBody(t *testing.T) {
	resetSharedClient(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "ok")
		_, _ = w.Write([]byte("payload"))
	}))
	defer server.Close()

	mapData, status, err := MakeCall(http.MethodGet, server.URL, nil, nil, nil, newTestSettings())
	if err != nil {
		t.Fatalf("MakeCall returned error: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("MakeCall returned status %d, want %d", status, http.StatusOK)
	}

	body, ok := mapData["body"].(string)
	if !ok {
		t.Fatalf("body entry has type %T, want string", mapData["body"])
	}
	if body != "payload" {
		t.Fatalf("body = %q, want %q", body, "payload")
	}

	header, ok := mapData["X-Test"].([]string)
	if !ok {
		t.Fatalf("X-Test header has type %T, want []string", mapData["X-Test"])
	}
	if len(header) != 1 || header[0] != "ok" {
		t.Fatalf("X-Test header = %v, want [ok]", header)
	}
}

func TestHTTPClientProcessTreats2xxAsSuccess(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
	}{
		{
			name:   "created",
			status: http.StatusCreated,
			body:   "created",
		},
		{
			name:   "no content",
			status: http.StatusNoContent,
			body:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSharedClient(t)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				if tt.body != "" {
					_, _ = w.Write([]byte(tt.body))
				}
			}))
			defer server.Close()

			httpClient := NewHTTPClient(http.MethodGet, server.URL, nil, nil, newTestSettings())

			mapData, status, err := httpClient.Process(nil)
			if err != nil {
				t.Fatalf("Process returned error: %v", err)
			}
			if status != tt.status {
				t.Fatalf("Process returned status %d, want %d", status, tt.status)
			}

			body, ok := mapData["body"].(string)
			if !ok {
				t.Fatalf("body entry has type %T, want string", mapData["body"])
			}
			if body != tt.body {
				t.Fatalf("body = %q, want %q", body, tt.body)
			}
		})
	}
}

func TestHTTPClientProcessReturnsStatusAndBodyOnErrorResponse(t *testing.T) {
	resetSharedClient(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	httpClient := NewHTTPClient(http.MethodGet, server.URL, nil, nil, newTestSettings())

	mapData, status, err := httpClient.Process(nil)
	if err == nil {
		t.Fatal("Process returned nil error, want non-nil")
	}
	if status != http.StatusInternalServerError {
		t.Fatalf("Process returned status %d, want %d", status, http.StatusInternalServerError)
	}

	body, ok := mapData["body"].(string)
	if !ok {
		t.Fatalf("body entry has type %T, want string", mapData["body"])
	}
	if strings.TrimSpace(body) != "boom" {
		t.Fatalf("body = %q, want %q", body, "boom")
	}
}

func TestMakeCallRawUsesBackgroundForNilContext(t *testing.T) {
	resetSharedClient(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("raw"))
	}))
	defer server.Close()

	resp, status, err := MakeCallRaw(nil, http.MethodGet, server.URL, nil, nil, nil, newTestSettings())
	if err != nil {
		t.Fatalf("MakeCallRaw returned error: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("MakeCallRaw returned status %d, want %d", status, http.StatusOK)
	}
	if resp == nil {
		t.Fatal("MakeCallRaw returned nil response")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}
	if string(body) != "raw" {
		t.Fatalf("body = %q, want %q", string(body), "raw")
	}
}

func TestMakeCallRawReturnsCatchAllStatusWhenContextCanceled(t *testing.T) {
	resetSharedClient(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("unexpected"))
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	resp, status, err := MakeCallRaw(ctx, http.MethodGet, server.URL, nil, nil, nil, newTestSettings())
	if err == nil {
		t.Fatal("MakeCallRaw returned nil error, want non-nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context canceled", err)
	}
	if resp != nil {
		t.Fatalf("MakeCallRaw returned response %#v, want nil", resp)
	}
	if status != 1000 {
		t.Fatalf("MakeCallRaw returned status %d, want %d", status, 1000)
	}
}
