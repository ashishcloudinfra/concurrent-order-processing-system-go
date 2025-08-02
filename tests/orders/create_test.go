package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"concurrent-order-processing-system/routes"
)

func TestCreateOrderHandler(t *testing.T) {
	router := routes.InitializeRoutes()

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		checkResponse  bool
	}{
		{
			name: "Valid Order",
			payload: map[string]interface{}{
				"customer": "John Doe",
				"items": []map[string]interface{}{
					{"name": "Product 1", "quantity": 2},
					{"name": "Product 2", "quantity": 1},
				},
			},
			expectedStatus: http.StatusCreated,
			checkResponse:  true,
		},
		{
			name: "Missing Customer",
			payload: map[string]interface{}{
				"items": []map[string]interface{}{
					{"name": "Product 1", "quantity": 2},
				},
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  false,
		},
		{
			name:           "Invalid JSON",
			payload:        nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			
			if tt.payload != nil {
				body, err = json.Marshal(tt.payload)
				if err != nil {
					t.Fatalf("Failed to marshal test payload: %v", err)
				}
			} else {
				body = []byte(`{invalid json`)
			}

			req, err := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.checkResponse {
				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Fatalf("Data field not found in response or not an object")
				}

				if _, ok := data["id"].(string); !ok {
					t.Errorf("ID field not found in response or not a string")
				}

				status, ok := data["status"].(string)
				if !ok || status != "pending" {
					t.Errorf("Status field not correct: got %v want pending", status)
				}
			}
		})
	}
}
