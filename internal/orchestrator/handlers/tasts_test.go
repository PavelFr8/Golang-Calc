package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PavelFr8/Golang-Calc/internal/orchestrator/handlers"
	"github.com/PavelFr8/Golang-Calc/internal/orchestrator/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupTaskRouter() *mux.Router {
	r := mux.NewRouter()
	handlers.RegisterTaskHandlers(r)
	return r
}

func TestHandleGetTask(t *testing.T) {
	models.TasksQueue = []models.Task{
		{ID: 1, Arg1: 2, Arg2: 3, Operation: "+"},
	}

	req, _ := http.NewRequest("GET", "/internal/task", nil)
	rec := httptest.NewRecorder()

	setupTaskRouter().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]models.Task
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp["task"].ID)
}

func TestHandleGetTask_NoTasks(t *testing.T) {
	models.TasksQueue = []models.Task{}

	req, _ := http.NewRequest("GET", "/internal/task", nil)
	rec := httptest.NewRecorder()

	setupTaskRouter().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
