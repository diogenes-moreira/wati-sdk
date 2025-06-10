package wati

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		token    string
		options  []ClientOption
		wantErr  bool
	}{
		{
			name:     "valid client creation",
			endpoint: "https://test.wati.io",
			token:    "test-token",
			options:  nil,
			wantErr:  false,
		},
		{
			name:     "client with options",
			endpoint: "https://test.wati.io",
			token:    "test-token",
			options:  []ClientOption{WithTimeout(30), WithRetries(3)},
			wantErr:  false,
		},
		{
			name:     "endpoint with trailing slash",
			endpoint: "https://test.wati.io/",
			token:    "test-token",
			options:  nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.endpoint, tt.token, tt.options...)
			
			if client == nil {
				t.Error("NewClient() returned nil")
				return
			}

			// Verificar que los servicios están inicializados
			if client.Contacts() == nil {
				t.Error("Contacts service not initialized")
			}
			
			if client.Messages() == nil {
				t.Error("Messages service not initialized")
			}
			
			if client.Chatbots() == nil {
				t.Error("Chatbots service not initialized")
			}
			
			if client.Media() == nil {
				t.Error("Media service not initialized")
			}
			
			if client.Webhooks() == nil {
				t.Error("Webhooks service not initialized")
			}
		})
	}
}

func TestClientOptions(t *testing.T) {
	client := NewClient(
		"https://test.wati.io",
		"test-token",
		WithTimeout(45),
		WithRetries(5),
		WithUserAgent("TestAgent/1.0"),
	)

	config := client.GetConfig()
	
	if config.Timeout != 45*time.Second {
		t.Errorf("Expected timeout 45s, got %v", config.Timeout)
	}
	
	if config.MaxRetries != 5 {
		t.Errorf("Expected max retries 5, got %d", config.MaxRetries)
	}
	
	if config.UserAgent != "TestAgent/1.0" {
		t.Errorf("Expected user agent 'TestAgent/1.0', got %s", config.UserAgent)
	}
}

func TestClientDoRequest(t *testing.T) {
	// Crear servidor de prueba
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar headers
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", r.Header.Get("Authorization"))
		}
		
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got %s", r.Header.Get("Content-Type"))
		}

		// Respuesta de prueba
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": true, "message": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	
	ctx := context.Background()
	
	var response struct {
		Result  bool   `json:"result"`
		Message string `json:"message"`
	}
	
	err := client.DoRequest(ctx, "GET", "/test", nil, &response)
	if err != nil {
		t.Errorf("DoRequest() error = %v", err)
		return
	}
	
	if !response.Result {
		t.Error("Expected result to be true")
	}
	
	if response.Message != "success" {
		t.Errorf("Expected message 'success', got %s", response.Message)
	}
}

func TestClientDoRequestWithError(t *testing.T) {
	// Servidor que retorna error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"result": false, "error": "invalid request"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	
	ctx := context.Background()
	
	var response struct {
		Result bool   `json:"result"`
		Error  string `json:"error"`
	}
	
	err := client.DoRequest(ctx, "POST", "/test", nil, &response)
	if err == nil {
		t.Error("Expected error but got nil")
		return
	}
	
	// Verificar que es un APIError
	if apiErr, ok := err.(*APIError); ok {
		if apiErr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, apiErr.Code)
		}
	} else {
		t.Errorf("Expected APIError, got %T", err)
	}
}

func TestClientRateLimit(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": true}`))
	}))
	defer server.Close()

	// Cliente con rate limit muy bajo para testing
	client := NewClient(server.URL, "test-token", WithRateLimit(2, 1)) // 2 requests por segundo
	
	ctx := context.Background()
	
	start := time.Now()
	
	// Hacer 3 requests
	for i := 0; i < 3; i++ {
		var response struct {
			Result bool `json:"result"`
		}
		
		err := client.DoRequest(ctx, "GET", "/test", nil, &response)
		if err != nil {
			t.Errorf("Request %d failed: %v", i+1, err)
		}
	}
	
	elapsed := time.Since(start)
	
	// Debería tomar al menos 1 segundo debido al rate limiting
	if elapsed < 500*time.Millisecond {
		t.Errorf("Rate limiting not working properly, elapsed: %v", elapsed)
	}
	
	if requestCount != 3 {
		t.Errorf("Expected 3 requests, got %d", requestCount)
	}
}

func TestClientTimeout(t *testing.T) {
	// Servidor que tarda mucho en responder
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", WithTimeout(1)) // 1 segundo timeout
	
	ctx := context.Background()
	
	var response struct{}
	
	start := time.Now()
	err := client.DoRequest(ctx, "GET", "/test", nil, &response)
	elapsed := time.Since(start)
	
	if err == nil {
		t.Error("Expected timeout error but got nil")
		return
	}
	
	// Debería fallar en aproximadamente 1 segundo
	if elapsed > 1500*time.Millisecond {
		t.Errorf("Timeout took too long: %v", elapsed)
	}
}

func TestClientRetries(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		
		// Fallar las primeras 2 requests, éxito en la tercera
		if requestCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", WithRetries(3))
	
	ctx := context.Background()
	
	var response struct {
		Result bool `json:"result"`
	}
	
	err := client.DoRequest(ctx, "GET", "/test", nil, &response)
	if err != nil {
		t.Errorf("Request failed after retries: %v", err)
		return
	}
	
	if requestCount != 3 {
		t.Errorf("Expected 3 requests (2 retries), got %d", requestCount)
	}
	
	if !response.Result {
		t.Error("Expected successful response")
	}
}

func TestClientContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	
	var response struct{}
	
	start := time.Now()
	err := client.DoRequest(ctx, "GET", "/test", nil, &response)
	elapsed := time.Since(start)
	
	if err == nil {
		t.Error("Expected context cancellation error but got nil")
		return
	}
	
	// Debería cancelarse en aproximadamente 500ms
	if elapsed > 1*time.Second {
		t.Errorf("Context cancellation took too long: %v", elapsed)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", config.Timeout)
	}
	
	if config.MaxRetries != 3 {
		t.Errorf("Expected default max retries 3, got %d", config.MaxRetries)
	}
	
	if config.UserAgent == "" {
		t.Error("Expected default user agent to be set")
	}
}

// Benchmark para medir performance
func BenchmarkClientDoRequest(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	ctx := context.Background()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		var response struct {
			Result bool `json:"result"`
		}
		
		client.DoRequest(ctx, "GET", "/test", nil, &response)
	}
}

