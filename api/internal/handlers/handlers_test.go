package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"promptforge/internal/models"
)

func TestHealthCheckEndpoint(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a minimal handlers instance for testing
	h := &Handlers{}

	err := h.HealthCheck(c)
	if err != nil {
		t.Fatalf("HealthCheck returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}

	if response["service"] != "PromptForge API" {
		t.Errorf("Expected service 'PromptForge API', got '%s'", response["service"])
	}
}

func TestRequestValidation(t *testing.T) {
	// Test JSON binding directly with Echo context
	tests := []struct {
		name        string
		requestBody string
		shouldError bool
	}{
		{
			name:        "Valid JSON",
			requestBody: `{"prompt": "Test prompt", "model": "gpt-4", "temperature": 0.7}`,
			shouldError: false,
		},
		{
			name:        "Invalid JSON",
			requestBody: `{"prompt": "Test prompt", "model": "gpt-4", "temperature":}`,
			shouldError: true,
		},
		{
			name:        "Empty JSON object",
			requestBody: `{}`,
			shouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(test.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			var executeReq models.ExecuteRequest
			err := c.Bind(&executeReq)

			if test.shouldError && err == nil {
				t.Error("Expected error for invalid JSON but got none")
			}
			if !test.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestModelStructures(t *testing.T) {
	// Test that our request/response structures serialize correctly
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "ExecuteRequest",
			data: models.ExecuteRequest{
				Prompt:      "Test prompt",
				Model:       "gpt-4",
				Temperature: 0.7,
				MaxTokens:   100,
			},
		},
		{
			name: "CritiqueRequest",
			data: models.CritiqueRequest{
				Prompt: "Test prompt for critique",
				Model:  "gpt-4",
			},
		},
		{
			name: "SaveHistoryRequest",
			data: models.SaveHistoryRequest{
				Prompt:      "Test prompt",
				Model:       "gpt-4",
				Temperature: 0.7,
				MaxTokens:   100,
				Success:     true,
				Response:    "Test response",
			},
		},
		{
			name: "APIResponse Success",
			data: models.APIResponse{
				Success: true,
				Data:    "Operation successful",
			},
		},
		{
			name: "APIResponse Error",
			data: models.APIResponse{
				Success: false,
				Error:   "Something went wrong",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Test marshaling
			jsonData, err := json.Marshal(test.data)
			if err != nil {
				t.Fatalf("Failed to marshal %s: %v", test.name, err)
			}

			// Test that we get valid JSON
			var result map[string]interface{}
			err = json.Unmarshal(jsonData, &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal %s as generic JSON: %v", test.name, err)
			}

			// Verify we have some data
			if len(result) == 0 {
				t.Errorf("Expected non-empty JSON for %s", test.name)
			}
		})
	}
}

func TestHTTPStatusCodes(t *testing.T) {
	// Test that handlers return appropriate HTTP status codes for different scenarios
	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
		handler        func(*Handlers, echo.Context) error
	}{
		{
			name:           "Health check returns 200",
			method:         "GET",
			path:           "/api/health",
			body:           "",
			expectedStatus: http.StatusOK,
			handler:        (*Handlers).HealthCheck,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			var req *http.Request
			if test.body != "" {
				req = httptest.NewRequest(test.method, test.path, bytes.NewReader([]byte(test.body)))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(test.method, test.path, nil)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &Handlers{}
			err := test.handler(h, c)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			if rec.Code != test.expectedStatus {
				t.Errorf("Expected status code %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}

func TestResponseFormats(t *testing.T) {
	// Test that responses follow the expected format
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &Handlers{}
	err := h.HealthCheck(c)
	if err != nil {
		t.Fatalf("HealthCheck returned error: %v", err)
	}

	// Verify response is valid JSON
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	// Verify response has expected fields
	if _, exists := response["status"]; !exists {
		t.Error("Response missing 'status' field")
	}

	if _, exists := response["service"]; !exists {
		t.Error("Response missing 'service' field")
	}
}

func TestJSONBinding(t *testing.T) {
	// Test that Echo's JSON binding works correctly with our models
	tests := []struct {
		name        string
		jsonString  string
		target      interface{}
		shouldError bool
	}{
		{
			name:        "Valid ExecuteRequest",
			jsonString:  `{"prompt": "test", "model": "gpt-4", "temperature": 0.7, "max_tokens": 100}`,
			target:      &models.ExecuteRequest{},
			shouldError: false,
		},
		{
			name:        "Invalid JSON",
			jsonString:  `{"prompt": "test", "temperature":}`,
			target:      &models.ExecuteRequest{},
			shouldError: true,
		},
		{
			name:        "Valid CritiqueRequest",
			jsonString:  `{"prompt": "test prompt", "model": "gpt-4"}`,
			target:      &models.CritiqueRequest{},
			shouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(test.jsonString)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := c.Bind(test.target)
			if test.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !test.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
