package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"concurrent-order-processing-system/utils/responses"
)

func TestWriteSuccessResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	testData := map[string]string{"test": "value"}
	
	responses.WriteSuccessResponse(rr, "Success message", testData)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	status, ok := response["status"].(string)
	if !ok || status != "success" {
		t.Errorf("Unexpected status in response: got %v", status)
	}
	
	message, ok := response["message"].(string)
	if !ok || message != "Success message" {
		t.Errorf("Unexpected message in response: got %v", message)
	}
	
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Data field not found in response or not an object")
	}
	
	testValue, ok := data["test"].(string)
	if !ok || testValue != "value" {
		t.Errorf("Unexpected data in response: got %v", testValue)
	}
}

func TestWriteSuccessResponseWithCode(t *testing.T) {
	rr := httptest.NewRecorder()
	testData := map[string]string{"test": "value"}
	
	responses.WriteSuccessResponseWithCode(rr, "Created message", testData, http.StatusCreated)
	
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	status, ok := response["status"].(string)
	if !ok || status != "success" {
		t.Errorf("Unexpected status in response: got %v", status)
	}
}

func TestWriteErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	
	responses.WriteErrorResponse(rr, "Error message", http.StatusBadRequest)
	
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	status, ok := response["status"].(string)
	if !ok || status != "error" {
		t.Errorf("Unexpected status in response: got %v", status)
	}
	
	message, ok := response["message"].(string)
	if !ok || message != "Error message" {
		t.Errorf("Unexpected message in response: got %v", message)
	}
}
