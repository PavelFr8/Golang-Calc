package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PavelFr8/Golang-Calc/internal/orchestrator/handlers"
	"github.com/PavelFr8/Golang-Calc/internal/orchestrator/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	handlers.RegisterExpressionHandlers(r)
	return r
}

func TestHandleCalculateExpression(t *testing.T) {
	models.Expressions = make(map[int]*models.Expression) // Очищаем хранилище

	reqBody := []byte(`{"expression": "3 + 5 * 2 * 2.5"}`)
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	log.Printf("Sending request to /api/v1/calculate with body: %s", reqBody)
	setupRouter().ServeHTTP(rec, req)

	log.Printf("Response Code: %d", rec.Code)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]int
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		log.Printf("Error unmarshalling response: %v", err)
	}
	assert.NoError(t, err)
	assert.Greater(t, resp["id"], 0)
	log.Printf("Received ID: %d", resp["id"])
}

func TestHandleGetExpressions(t *testing.T) {
	models.Expressions = make(map[int]*models.Expression)
	models.Expressions[1] = &models.Expression{ID: 1, Status: "pending"}

	req, _ := http.NewRequest("GET", "/api/v1/expressions", nil)
	rec := httptest.NewRecorder()

	log.Println("Sending request to /api/v1/expressions")
	setupRouter().ServeHTTP(rec, req)

	log.Printf("Response Code: %d", rec.Code)
	assert.Equal(t, http.StatusOK, rec.Code)
	log.Printf("Response Body: %s", rec.Body.String())
	assert.Contains(t, rec.Body.String(), `"id":1`)
}

func TestHandleGetExpressionByID(t *testing.T) {
	models.Expressions = make(map[int]*models.Expression)
	models.Expressions[2] = &models.Expression{ID: 2, Status: "completed"}

	req, _ := http.NewRequest("GET", "/api/v1/expressions/2", nil)
	rec := httptest.NewRecorder()

	log.Println("Sending request to /api/v1/expressions/2")
	setupRouter().ServeHTTP(rec, req)

	log.Printf("Response Code: %d", rec.Code)
	assert.Equal(t, http.StatusOK, rec.Code)
	log.Printf("Response Body: %s", rec.Body.String())
	assert.Contains(t, rec.Body.String(), `"id":2`)
}

func TestHandleGetExpressionByID_NotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/expressions/99", nil)
	rec := httptest.NewRecorder()

	log.Println("Sending request to /api/v1/expressions/99")
	setupRouter().ServeHTTP(rec, req)

	log.Printf("Response Code: %d", rec.Code)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
