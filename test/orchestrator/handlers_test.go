package orchestrator_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PavelFr8/Golang-Calc/internal/orchestrator"
	"github.com/gorilla/mux"
)

type registerPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type exprPayload struct {
	Expression string `json:"expression"`
}

func SetupRouterForTest(o *orchestrator.Orchestrator) *mux.Router {
	r := mux.NewRouter()
	exempt := map[string]bool{
		"/api/v1/login":    true,
		"/api/v1/register": true,
		"/":                true,
	}

	r.Use(orchestrator.JWTMiddleware(o.Config.JWTsecret, exempt))

	r.HandleFunc("/api/v1/calculate", o.CalculateHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/expressions", o.ExpressionsHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/expressions/{id}", o.ExpressionByIDHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/register", o.RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/login", o.LoginHandler).Methods(http.MethodPost)
	return r
}

func TestFullFlow_RegisterLoginCalculate(t *testing.T) {
	o := orchestrator.New()
	router := SetupRouterForTest(o)
	defer func() {
		sqlDB, _ := o.R.DB.DB()
		sqlDB.Close()
		os.Remove("./database.db")
	}()
	// --- Регистрация ---
	reg := registerPayload{Login: "testuser", Password: "testpass"}
	body, _ := json.Marshal(reg)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusContinue {
		t.Fatalf("Регистрация не удалась: %v", resp.Body.String())
	}

	// --- Логинизация ---
	req = httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusContinue {
		t.Fatalf("Логинизация не удалась: %v", resp.Body.String())
	}

	var loginResp loginResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &loginResp)
	if loginResp.Token == "" {
		t.Fatal("JWT токен не получен")
	}

	// --- Отправка выражения ---
	expr := exprPayload{Expression: "2 + 2"}
	body, _ = json.Marshal(expr)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	t.Log("Токен: ", "Bearer "+loginResp.Token)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("Ошибка добавления выражения: %v", resp.Body.String())
	}
	
	
}
