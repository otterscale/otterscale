package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPGet_Success(t *testing.T) {
	// Create a test server that returns a successful response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	ctx := context.Background()
	result, err := HTTPGet(ctx, server.URL)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := "test response"
	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, string(result))
	}
}

func TestHTTPGet_ContextCancellation(t *testing.T) {
	// Create a test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("delayed response"))
	}))
	defer server.Close()

	// Create a context that cancels quickly
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := HTTPGet(ctx, server.URL)

	if err == nil {
		t.Error("Expected an error due to context cancellation, got nil")
	}
}

func TestHTTPGet_HTTPErrors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedErrMsg string
	}{
		{
			name:           "404 Not Found",
			statusCode:     http.StatusNotFound,
			expectedErrMsg: "failed with code 404",
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: "failed with code 500",
		},
		{
			name:           "400 Bad Request",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: "failed with code 400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			ctx := context.Background()
			_, err := HTTPGet(ctx, server.URL)

			if err == nil {
				t.Errorf("Expected an error for status code %d, got nil", tt.statusCode)
			}

			if err != nil && !contains(err.Error(), tt.expectedErrMsg) {
				t.Errorf("Expected error to contain %q, got %q", tt.expectedErrMsg, err.Error())
			}
		})
	}
}

func TestHTTPGet_SuccessfulStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"200 OK", http.StatusOK},
		{"201 Created", http.StatusCreated},
		{"202 Accepted", http.StatusAccepted},
		{"204 No Content", http.StatusNoContent},
		{"299 Custom Success", 299},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte("success"))
			}))
			defer server.Close()

			ctx := context.Background()
			result, err := HTTPGet(ctx, server.URL)

			if err != nil {
				t.Errorf("Expected no error for status code %d, got %v", tt.statusCode, err)
			}

			if tt.statusCode != http.StatusNoContent && string(result) != "success" {
				t.Errorf("Expected 'success', got %q", string(result))
			}
		})
	}
}

func TestHTTPGet_InvalidURL(t *testing.T) {
	ctx := context.Background()
	_, err := HTTPGet(ctx, "invalid-url")

	if err == nil {
		t.Error("Expected an error for invalid URL, got nil")
	}
}

func TestHTTPGet_NetworkError(t *testing.T) {
	ctx := context.Background()
	// Use a URL that should cause a network error
	_, err := HTTPGet(ctx, "http://localhost:99999/nonexistent")

	if err == nil {
		t.Error("Expected a network error, got nil")
	}
}

func TestHTTPGet_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Don't write any content
	}))
	defer server.Close()

	ctx := context.Background()
	result, err := HTTPGet(ctx, server.URL)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty response, got %d bytes", len(result))
	}
}

func TestHTTPGet_LargeResponse(t *testing.T) {
	largeContent := make([]byte, 1024*1024) // 1MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(largeContent)
	}))
	defer server.Close()

	ctx := context.Background()
	result, err := HTTPGet(ctx, server.URL)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != len(largeContent) {
		t.Errorf("Expected %d bytes, got %d bytes", len(largeContent), len(result))
	}
}

func TestHTTPGet_JSONResponse(t *testing.T) {
	jsonResponse := `{"message": "hello", "status": "ok"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
	defer server.Close()

	ctx := context.Background()
	result, err := HTTPGet(ctx, server.URL)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(result) != jsonResponse {
		t.Errorf("Expected %q, got %q", jsonResponse, string(result))
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && s[:len(substr)] == substr) ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkHTTPGet(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("benchmark response"))
	}))
	defer server.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HTTPGet(ctx, server.URL)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkHTTPGet_Parallel(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("parallel benchmark response"))
	}))
	defer server.Close()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := HTTPGet(ctx, server.URL)
			if err != nil {
				b.Fatalf("Unexpected error: %v", err)
			}
		}
	})
}
