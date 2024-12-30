package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PavelFr8/Golang-Calc/internal/application"
)

func TestRequestHandler(t *testing.T) {
	tt := []struct {
		method         string
		expression     string
		expectedStatus int
		expectedError  string
		expectedResult string
	}{
		// Valid expression
		{
			method:         http.MethodPost,
			expression:     `{"expression": "20 + 2 * ((7*7) / 7)"}`,
			expectedStatus: http.StatusOK,
			expectedResult: "34.000000",
		},
		// Invalid expression (unprocessable entity)
		{
			method:         http.MethodPost,
			expression:     `{"expression": "invalid_expression"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "Expression is not valid",
		},
		// Divide by zero (internal server error)
		{
			method:         http.MethodPost,
			expression:     `{"expression": "7777/0"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Internal server error",
		},
		// Correct method, but invalid HTTP method (Method Not Allowed)
		{
			method:         http.MethodPut,
			expression:     `{"expression": "2+2"}`,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, test_case := range tt {
		req := httptest.NewRequest(test_case.method, "/api/v1/calculate", bytes.NewBufferString(test_case.expression))
		w := httptest.NewRecorder()

		application.CalcHandler()(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != test_case.expectedStatus {
			t.Errorf("expected status code %d, got %d", test_case.expectedStatus, res.StatusCode)
		}

		if test_case.expectedError != "" {
			var response application.Response
			if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
				t.Errorf("failed to decode response: %v", err)
			}
			if response.Error != test_case.expectedError {
				t.Errorf("expected error message %q, got %q", test_case.expectedError, response.Error)
			}
		}

		if test_case.expectedResult != "" {
			var response application.Response
			if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
				t.Errorf("failed to decode response: %v", err)
			}
			if response.Result != test_case.expectedResult {
				t.Errorf("expected result %q, got %q", test_case.expectedResult, response.Result)
			}
		}
	}
}
