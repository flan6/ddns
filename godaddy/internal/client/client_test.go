package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flan6/ddns/godaddy/entity"
)

func TestGoDaddyClient_Request(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		case "/error":
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.GoDaddyError{Code: "BadRequest", Message: "Invalid request"})
		default:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(entity.GoDaddyError{Code: "NotFound", Message: "Resource not found"})
		}
	}))
	defer server.Close()

	// Test cases
	tests := []struct {
		name     string
		method   string
		url      string
		data     []byte
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "Successful request",
			method:  "GET",
			url:     server.URL + "/success",
			wantErr: false,
		},
		{
			name:     "GoDaddy error",
			method:   "GET",
			url:      server.URL + "/error",
			wantErr:  true,
			errorMsg: "BadRequest: Invalid request",
		},
	}

	// Initialize client
	client := GoDaddyClient{
		httpClient: server.Client(),
		secret:     entity.Secret{ApiKey: "apikey", ApiSecret: "apisecret"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Request(tt.method, tt.url, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && err.Error() != tt.errorMsg {
				t.Errorf("Request() error = %v, errorMsg %v", err, tt.errorMsg)
			}
		})
	}
}
