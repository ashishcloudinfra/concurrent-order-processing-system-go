package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"concurrent-order-processing-system/routes"
)

func TestGetOrderStatusHandler(t *testing.T) {
	router := routes.InitializeRoutes()

	createOrderPayload := map[string]interface{}{
		"customer": "Test Customer",
		"items": []map[string]interface{}{
			{"name": "Test Item", "quantity": 1},
		},
	}

	body, _ := json.Marshal(createOrderPayload)
	createReq, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	var createResp map[string]interface{}
	json.NewDecoder(createRec.Body).Decode(&createResp)
	orderData := createResp["data"].(map[string]interface{})
	orderID := orderData["id"].(string)

	tests := []struct {
		name           string
		orderID        string
		expectedStatus int
		expectedValue  string
	}{
		{
			name:           "Valid Order ID",
			orderID:        orderID,
			expectedStatus: http.StatusOK,
			expectedValue:  "pending",
		},
		{
			name:           "Invalid Order ID",
			orderID:        "nonexistent-id",
			expectedStatus: http.StatusNotFound,
			expectedValue:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/orders/status/%s", tt.orderID), nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Fatalf("Data field not found in response or not an object")
				}

				if status, ok := data["status"].(string); !ok || status != tt.expectedValue {
					t.Errorf("Status field not correct: got %v want %v", status, tt.expectedValue)
				}
			}
		})
	}
}
