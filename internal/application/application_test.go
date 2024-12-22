package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PavelFr8/Golang-Calc/internal/application"
)

func TestRequestHandlerBadMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/v1/calculate", bytes.NewBufferString(`{"expression": "2+2"}`))
	w := httptest.NewRecorder()
	application.CalcHandler(nil)(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, res.StatusCode)
	}
}

func TestRequestHandlerInvalidExpression(t *testing.T) {
	tt := []struct {
		expression     string
		expectedStatus int
		expectedError  string
	}{
		{`{"expression": "7777/0"}`, http.StatusInternalServerError, "Internal server error"},
		{`{"expression": "invalid_expression"}`, http.StatusUnprocessableEntity, "Expression is not valid"},
	}

	for _, test_case := range tt {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(test_case.expression))
		w := httptest.NewRecorder()

		application.CalcHandler(nil)(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != test_case.expectedStatus {
			t.Errorf("expected status code %d, got %d", test_case.expectedStatus, res.StatusCode)
		}

		var response application.Response
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}
		if response.Error != test_case.expectedError {
			t.Errorf("expected error message %q, got %q", test_case.expectedError, response.Error)
		}
	}
}

func TestRequestHandlerInternalServerError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(`{"expression": "1/0"}`))
	w := httptest.NewRecorder()

	application.CalcHandler(nil)(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, res.StatusCode)
	}

	var response application.Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	if response.Error != "Internal server error" {
		t.Errorf("expected error message %q, got %q", "internal server error", response.Error)
	}
}
