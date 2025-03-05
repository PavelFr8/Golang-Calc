package orchestrator

import (
	"net/http"

	"github.com/PavelFr8/Golang-Calc/internal/orchestrator/handlers"
	"github.com/PavelFr8/Golang-Calc/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Orchestrator struct {
	config *OrchestratorConfig
	logger *zap.Logger
}

func New() *Orchestrator {
	return &Orchestrator{
		config: NewOrchestratorConfig(),
		logger: logger.SetupLogger(),
	}
}

func (o *Orchestrator) RunServer() error {
	r := mux.NewRouter()

	// Добавляем мидлварь для логирования
	r.Use(logger.LoggingMiddleware(o.logger)) 

	// Регистрация маршрутов
	handlers.RegisterExpressionHandlers(r)
	handlers.RegisterTaskHandlers(r)

	o.logger.Info("Оркестратор-Сервер запущен", zap.String("address", ":"+o.config.OrchestratorPort))

	return http.ListenAndServe(":"+o.config.OrchestratorPort, r)
}
