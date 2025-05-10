package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/PavelFr8/Golang-Calc/pkg/tree"
)

func (o *Orchestrator) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Expression == "" {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusUnprocessableEntity)
		return
	}
	tree, err := tree.BuildNode(req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnprocessableEntity)
		return
	}
	o.Mu.Lock()
	defer o.Mu.Unlock()
	exprID := o.r.GetMaxExpressionID() + 1
	expr := &Expression{
		Expr:   req.Expression,
		Status: "pending",
		Node:    tree,
	}
	if expr.Node.IsLeaf {
		expr.Status = "completed"
	}
	o.r.CreateExpression(expr)
	o.Expressions[exprID] = expr
	o.NewTask(expr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uint{"id": exprID})
}

func (o *Orchestrator) ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	exprs := make([]*Expression, 0, len(o.Expressions))
	for _, expr := range o.Expressions {
		if expr.Node != nil && expr.Node.IsLeaf {
			if err := tree.Check(expr.Node); err != nil {
				expr.Result = nil
			} else {
				expr.Result = &expr.Node.Value
			}
		}
		exprs = append(exprs, expr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": exprs})
}

func (o *Orchestrator) ExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	int_id, _ := strconv.Atoi(r.URL.Path[len("/api/v1/expressions/"):])
	id := uint(int_id)
	o.Mu.Lock()
	expr, ok := o.Expressions[id]
	o.Mu.Unlock()
	if !ok {
		http.Error(w, `{"error":"Expression not found"}`, http.StatusNotFound)
		return
	}
	if expr.Node != nil && expr.Node.IsLeaf {
		expr.Result = &expr.Node.Value
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}

// Тупо отдаем самый последний элемент очереди
func (o *Orchestrator) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if len(o.TaskQueue) == 0 {
		http.Error(w, `{"error":"No task available"}`, http.StatusNotFound)
		return
	}

	task := o.TaskQueue[0]
	o.TaskQueue = o.TaskQueue[1:]

	if _, exists := o.Expressions[task.ExprID]; !exists {
		http.Error(w, `{"error":"Task expression not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"task": task})
}

// Тут уже с огромной болью с слезами, добавляем решенное выражение обратно в очередь
func (o *Orchestrator) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     uint  `json:"id"`
		Result float64 `json:"result"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.ID == 0 {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusUnprocessableEntity)
		return
	}
	o.Mu.Lock()
	task, ok := o.Tasks[req.ID]
	if !ok {
		o.Mu.Unlock()
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	}
	task.Node.IsLeaf = true
	task.Node.Value = req.Result
	delete(o.Tasks, req.ID)
	if expr, exists := o.Expressions[task.ExprID]; exists {
		o.NewTask(expr)
		if expr.Node.IsLeaf {
			expr.Status = "completed"
			expr.Result = &expr.Node.Value
			o.r.db.Updates(expr)
		}
	}
	o.Mu.Unlock()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"result accepted"}`))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join("web", "index.html")
	file, err := os.Open(indexPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "text/html")
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Unable to get file info", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, "index.html", fileInfo.ModTime(), file)
}