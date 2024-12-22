package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PavelFr8/Golang-Calc/pkg/calculation"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
	logger *zap.Logger
}

func setupLogger() *zap.Logger {
	// Настраиваем конфигурацию логгера
	config := zap.NewProductionConfig()
	// Уровень логирования
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Ошибка настройки логгера: %v\n", err)
		os.Exit(1) // Завершение программы
	}
	return logger
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
		logger: setupLogger(),
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,,omitempty"`
	Error  string `json:"error,omitempty"`
}

func loggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	// Middleware для логирования запросов и ответов
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)

			// Завершаем логирование после того, как запрос выполнен
			duration := time.Since(start)
			logger.Info("HTTP запрос",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		})
	}
}

func CalcHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Принимаем только POST!
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		request := new(Request)
		response := Response{}
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		result, err := calculation.Calc(request.Expression)
		if err != nil {
			if errors.Is(err, calculation.ErrInvalidExpression) {
				response.Error = err.Error()
				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(response)
				return
			} else if errors.Is(err, calculation.ErrCalculation) || errors.Is(err, calculation.ErrDivisionByZero) {
				response.Error = "Internal server error"
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			} else {
				response.Error = "Internal server error"
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		response.Result = fmt.Sprintf("%f", result)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func (a *Application) RunServer() error {
	r := mux.NewRouter()
	r.Use(loggingMiddleware(a.logger)) // Добавляем middleware для логирования

	// Принимаем только POST запросы!
	r.HandleFunc("/api/v1/calculate", CalcHandler(a.logger)).Methods("POST")

	a.logger.Info("Сервер запущен", zap.String("address", fmt.Sprintf(":%s", a.config.Addr))) // Логируем запуск сервера
	return http.ListenAndServe(":"+a.config.Addr, r)
}
