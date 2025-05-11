package orchestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/PavelFr8/Golang-Calc/pkg/hash"
	"github.com/PavelFr8/Golang-Calc/pkg/tree"
	"github.com/golang-jwt/jwt/v5"
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
	userID, ok := GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Auth fail. Refresh token"}`, http.StatusUnauthorized)
		return
	}
	expr := &Expression{
		Expr:   req.Expression,
		Status: "pending",
		Node:    tree,
		UserID: userID,
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

	userID, ok := GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Auth fail. Refresh token"}`, http.StatusUnauthorized)
		return
	}

	exprs := make([]*Expression, 0, len(o.Expressions))
	for _, expr := range o.Expressions {
		if expr.UserID == userID {
			if expr.Node != nil && expr.Node.IsLeaf {
				if err := tree.Check(expr.Node); err != nil {
					expr.Result = nil
				} else {
					expr.Result = expr.Node.Value
				}
			}
			exprs = append(exprs, expr)
		}
	}
	sort.Slice(exprs, func(i, j int) bool {return exprs[i].ID < exprs[j].ID})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": exprs})
}

func (o *Orchestrator) ExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r)
	if !ok {
		http.Error(w, `{"error":"Auth fail. Refresh token"}`, http.StatusUnauthorized)
		return
	}
	int_id, _ := strconv.Atoi(r.URL.Path[len("/api/v1/expressions/"):])
	id := uint(int_id)
	o.Mu.Lock()
	expr, ok := o.Expressions[id]
	if expr.UserID != userID {
		o.Mu.Unlock()
		http.Error(w, `{"error":"You haven't got access to this expression"}`, http.StatusForbidden)
		return
	}
	o.Mu.Unlock()
	if !ok {
		http.Error(w, `{"error":"Expression not found"}`, http.StatusNotFound)
		return
	}
	if expr.Node != nil && expr.Node.IsLeaf {
		expr.Result = expr.Node.Value
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
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

func (o *Orchestrator) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login      string   `json:"login"`
		Password   string   `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusUnprocessableEntity)
		return
	}
	err2 := o.r.CreateUser(req.Login, req.Password)
	if err2 != nil || req.Login == "" || req.Password == "" {
		http.Error(w, `{"error":"Invalid login or password"}`, http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusContinue)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func (o *Orchestrator) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login      string   `json:"login"`
		Password   string   `json:"password"`
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusUnprocessableEntity)
		return
	}
	err2 := o.r.db.Where("login = ?", req.Login).First(&user).Error
	if err2 != nil {
		http.Error(w, `{"error":"Invalid login or password"}`, http.StatusUnprocessableEntity)
		return
	}
	if err3 := hash.Compare(user.Password, req.Password); err3 != nil {
		http.Error(w, `{"error":"Invalid login or password"}`, http.StatusUnprocessableEntity)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})
	
	tokenString, err := token.SignedString(o.Config.JWTsecret)
	if err != nil {
		http.Error(w, `{"error":"Token generation failed"}`, http.StatusInternalServerError)
		fmt.Println("Error generating token:", err)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusContinue)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "OK",
		"token": tokenString})
}