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

func TestCancelOrderHandler(t *testing.T) {
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
	}{
		{
			name:           "Valid Cancel",
			orderID:        orderID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Order ID",
			orderID:        "nonexistent-id",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/orders/cancel/%s", tt.orderID), nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				getReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/orders/%s", tt.orderID), nil)
				getRr := httptest.NewRecorder()
				router.ServeHTTP(getRr, getReq)

				if getRr.Code != http.StatusOK {
					t.Errorf("Failed to get order after cancellation")
				}

				var getResp map[string]interface{}
				json.NewDecoder(getRr.Body).Decode(&getResp)
				updatedOrderData := getResp["data"].(map[string]interface{})
				updatedStatus := updatedOrderData["status"].(string)

				if updatedStatus != "canceled" {
					t.Errorf("Order status was not updated to canceled, got: %s", updatedStatus)
				}
			}
		})
	}
}
