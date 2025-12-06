package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"message": "test"}

	WriteJSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "test" {
		t.Errorf("Expected message 'test', got %s", response["message"])
	}
}

func TestWriteJSON_ComplexData(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"id":   123,
		"name": "test",
		"tags": []string{"tag1", "tag2"},
	}

	WriteJSON(w, http.StatusCreated, data)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["id"].(float64) != 123 {
		t.Errorf("Expected id 123, got %v", response["id"])
	}
}

func TestWriteJSON_ArrayData(t *testing.T) {
	w := httptest.NewRecorder()
	data := []string{"item1", "item2", "item3"}

	WriteJSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("Expected 3 items, got %d", len(response))
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	errorMsg := "something went wrong"

	WriteError(w, http.StatusBadRequest, errorMsg)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != errorMsg {
		t.Errorf("Expected error message '%s', got %s", errorMsg, response["error"])
	}
}

func TestWriteError_InternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	errorMsg := "internal server error"

	WriteError(w, http.StatusInternalServerError, errorMsg)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != errorMsg {
		t.Errorf("Expected error message '%s', got %s", errorMsg, response["error"])
	}
}

func TestWriteError_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	errorMsg := "resource not found"

	WriteError(w, http.StatusNotFound, errorMsg)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != errorMsg {
		t.Errorf("Expected error message '%s', got %s", errorMsg, response["error"])
	}
}
