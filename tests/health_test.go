package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"concurrent-order-processing-system/routes"
)

func TestHealthCheckHandler(t *testing.T) {
	router := routes.InitializeRoutes()

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	message, ok := response["message"].(string)
	if !ok || message != "Service is running" {
		t.Errorf("Unexpected message in response: got %v", message)
	}
}
