package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDelHopHeaders(t *testing.T) {
	header := http.Header{}
	header.Set("Connection", "keep-alive")
	header.Set("Keep-Alive", "timeout=5")
	header.Set("Proxy-Authenticate", "Basic")
	header.Set("Proxy-Authorization", "Basic abc123")
	header.Set("Te", "trailers")
	header.Set("Trailers", "X-Custom")
	header.Set("Transfer-Encoding", "chunked")
	header.Set("Upgrade", "websocket")
	header.Set("Content-Type", "application/json") // Should not be deleted

	delHopHeaders(header)

	for _, h := range hopHeaders {
		if header.Get(h) != "" {
			t.Errorf("hop header %q should have been deleted", h)
		}
	}

	if header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type header should not have been deleted")
	}
}

func TestCopyHeader(t *testing.T) {
	src := http.Header{}
	src.Set("Content-Type", "application/json")
	src.Add("X-Custom", "value1")
	src.Add("X-Custom", "value2")

	dst := http.Header{}
	copyHeader(dst, src)

	if dst.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type to be 'application/json', got %q", dst.Get("Content-Type"))
	}

	values := dst.Values("X-Custom")
	if len(values) != 2 {
		t.Errorf("expected 2 X-Custom values, got %d", len(values))
	}
}

func TestAppendHostToXForwardHeader(t *testing.T) {
	tests := []struct {
		name     string
		prior    []string
		host     string
		expected string
	}{
		{
			name:     "no prior X-Forwarded-For",
			prior:    nil,
			host:     "192.168.1.1",
			expected: "192.168.1.1",
		},
		{
			name:     "with single prior X-Forwarded-For",
			prior:    []string{"10.0.0.1"},
			host:     "192.168.1.1",
			expected: "10.0.0.1, 192.168.1.1",
		},
		{
			name:     "with multiple prior X-Forwarded-For",
			prior:    []string{"10.0.0.1", "10.0.0.2"},
			host:     "192.168.1.1",
			expected: "10.0.0.1, 10.0.0.2, 192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			for _, p := range tt.prior {
				header.Add("X-Forwarded-For", p)
			}

			appendHostToXForwardHeader(header, tt.host)

			got := header.Get("X-Forwarded-For")
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestProxyHTTP(t *testing.T) {
	// Create a test backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend", "test")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend response"))
	}))
	defer backend.Close()

	// Create a request to the backend
	req := httptest.NewRequest(http.MethodGet, backend.URL, nil)
	req.RemoteAddr = "127.0.0.1:12345"

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use the default transport for testing
	proxyHTTP(rr, req, http.DefaultTransport)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if rr.Header().Get("X-Backend") != "test" {
		t.Errorf("expected X-Backend header to be 'test', got %q", rr.Header().Get("X-Backend"))
	}

	if rr.Body.String() != "backend response" {
		t.Errorf("expected body 'backend response', got %q", rr.Body.String())
	}
}
