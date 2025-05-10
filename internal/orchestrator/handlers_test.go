package orchestrator_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PavelFr8/Golang-Calc/internal/orchestrator"
	"github.com/PavelFr8/Golang-Calc/pkg/logger"
	"github.com/PavelFr8/Golang-Calc/pkg/tree"
	"github.com/stretchr/testify/assert"
)

func TestNewOrchestrator(t *testing.T) {
	orchestrator := orchestrator.New()
	assert.NotNil(t, orchestrator)
	assert.NotNil(t, orchestrator.Expressions)
	assert.NotNil(t, orchestrator.Tasks)
	assert.NotNil(t, orchestrator.TaskQueue)
	assert.NotNil(t, orchestrator.Config)
}

func TestCalculateHandler(t *testing.T) {
	orchestrator := orchestrator.New()
	reqBody := []byte(`{"expression": "3+5"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	http.HandlerFunc(orchestrator.CalculateHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestExpressionsHandler(t *testing.T) {
	orchestrator := orchestrator.New()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/expressions", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	http.HandlerFunc(orchestrator.ExpressionsHandler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestExpressionByIDHandler_NotFound(t *testing.T) {
	orchestrator := orchestrator.New()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/expressions/999", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	http.HandlerFunc(orchestrator.ExpressionByIDHandler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestNewTask(t *testing.T) {
	o := orchestrator.New()
	logger.SetupLogger()
	node, _ := tree.BuildNode("3+5")
	expr := orchestrator.Expression{ID: uint(1), Expr: "3+5", Status: "pending", Node: node}
	o.NewTask(&expr)

	assert.NotEmpty(t, o.TaskQueue)
	assert.NotEmpty(t, o.Tasks)
}

func TestGetTaskHandler_NoTasks(t *testing.T) {
	orchestrator := orchestrator.New()
	req, _ := http.NewRequest(http.MethodGet, "/internal/task", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(orchestrator.GetTaskHandler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPostTaskHandler_TaskNotFound(t *testing.T) {
	orchestrator := orchestrator.New()
	reqBody := []byte(`{"id": "999", "result": 8}`)
	req, _ := http.NewRequest(http.MethodPost, "/internal/task", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	http.HandlerFunc(orchestrator.PostTaskHandler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
