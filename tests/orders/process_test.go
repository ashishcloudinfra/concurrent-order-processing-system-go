package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"concurrent-order-processing-system/routes"
)

func TestProcessOrdersHandler(t *testing.T) {
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

	req, _ := http.NewRequest(http.MethodPost, "/orders/process", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	time.Sleep(4 * time.Second)

	getReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/orders/%s", orderID), nil)
	getRr := httptest.NewRecorder()
	router.ServeHTTP(getRr, getReq)

	if getRr.Code != http.StatusOK {
		t.Errorf("Failed to get order after processing")
	}

	var getResp map[string]interface{}
	json.NewDecoder(getRr.Body).Decode(&getResp)
	updatedOrderData := getResp["data"].(map[string]interface{})
	updatedStatus := updatedOrderData["status"].(string)

	if updatedStatus == "pending" {
		t.Errorf("Order status was not updated after processing")
	}

	if updatedStatus != "shipped" && updatedStatus != "canceled" {
		t.Errorf("Order status has unexpected value: %s", updatedStatus)
	}
}

func TestProcessOrdersHandlerNoOrders(t *testing.T) {
	router := routes.InitializeRoutes()

	req, _ := http.NewRequest(http.MethodPost, "/orders/process", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
