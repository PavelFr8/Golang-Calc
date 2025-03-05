package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/PavelFr8/Golang-Calc/internal/orchestrator/models"
	"github.com/gorilla/mux"
)

func RegisterTaskHandlers(r *mux.Router) {
	r.HandleFunc("/internal/task", HandleGetTask).Methods("GET")
	r.HandleFunc("/internal/task", HandlePostTaskResult).Methods("POST")
}

// Тупо отдаем самый последний элемент очереди
func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	models.Mutex.Lock()
	defer models.Mutex.Unlock()

	if len(models.TasksQueue) == 0 {
		http.Error(w, "No task available", http.StatusNotFound)
		return
	}

	task := models.TasksQueue[0]
	models.TasksQueue = models.TasksQueue[1:]

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
}

// Тут уже с огромной болью с слезами, добавляем решенное выражение обратно в очередь
func HandlePostTaskResult(w http.ResponseWriter, r *http.Request) {
	var result struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	models.Mutex.Lock()
	defer models.Mutex.Unlock()

	for _, expr := range models.Expressions {
		for i, task := range expr.Tasks {
			if task.ID == result.ID {
				expr.Tasks = append(expr.Tasks[:i], expr.Tasks[i+1:]...)
				if len(expr.Tasks) == 0 {
					expr.Status = "completed"
					expr.Result = &result.Result
				}
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}