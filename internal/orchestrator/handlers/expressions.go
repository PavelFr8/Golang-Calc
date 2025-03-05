package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/PavelFr8/Golang-Calc/internal/orchestrator/models"
	"github.com/gorilla/mux"
)

func HandleCalculateExpression(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
		return
	}

	// Блокируем, чтобы корректно рассчитать ID
	models.Mutex.Lock()
	models.ExprCounter++
	exprID := models.ExprCounter
	models.Mutex.Unlock()
	
	// Парсим выражение
	expr, err := models.ParseExpression(exprID, request.Expression)
	if err != nil {
		http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
		return
	}

	// От греха подальше снова блокируем доступ, чтобы все аккуратненько добавить в очередь
	models.Mutex.Lock()
	models.Expressions[expr.ID] = expr
	models.TasksQueue = append(models.TasksQueue, expr.Tasks...)
	models.Mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": expr.ID})
}

// Ниже логика возврата всех expressions
func HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	models.Mutex.Lock()
	defer models.Mutex.Unlock()

	var response struct {
		Expressions []models.Expression `json:"expressions"`
	}

	for _, expr := range models.Expressions {
		response.Expressions = append(response.Expressions, *expr)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// returning expression by id logic below :)
func HandleGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	models.Mutex.Lock()
	defer models.Mutex.Unlock()

	expr, found := models.Expressions[id]
	if !found {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expr)
}


// Func for returning all handlers, which is connected with expressions
func RegisterExpressionHandlers(r *mux.Router) {
	r.HandleFunc("/api/v1/calculate", HandleCalculateExpression).Methods("POST")
	r.HandleFunc("/api/v1/expressions", HandleGetExpressions).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id:[0-9]+}", HandleGetExpressionByID).Methods("GET")
}
