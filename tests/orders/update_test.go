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

func TestUpdateOrderStatusHandler(t *testing.T) {
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
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:    "Valid Update",
			orderID: orderID,
			payload: map[string]interface{}{
				"status": "shipped",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Invalid Order ID",
			orderID: "nonexistent-id",
			payload: map[string]interface{}{
				"status": "shipped",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:    "Invalid Status",
			orderID: orderID,
			payload: map[string]interface{}{
				"status": "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateBody, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/orders/%s", tt.orderID), bytes.NewBuffer(updateBody))
			req.Header.Set("Content-Type", "application/json")
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

				if status, ok := data["status"].(string); !ok || status != tt.payload["status"] {
					t.Errorf("Status field not updated correctly: got %v want %v", status, tt.payload["status"])
				}
			}
		})
	}
}
