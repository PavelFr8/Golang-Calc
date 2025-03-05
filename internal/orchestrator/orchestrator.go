package orchestrator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PavelFr8/Golang-Calc/pkg/env"
	"github.com/PavelFr8/Golang-Calc/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type OrchestratorConfig struct {
	OrchestratorPort string
	TimeAddition time.Duration
	TimeSubtraction time.Duration
	TimeMultiplication time.Duration
	TimeDivision time.Duration
}

func NewOrchestratorConfig() *OrchestratorConfig {
	return &OrchestratorConfig{
		OrchestratorPort: env.GetEnv("ORCHESTRATOR_PORT", "8081"),
		TimeAddition: time.Duration(env.GetEnvAsInt("TIME_ADDITION", 1000)),
		TimeSubtraction: time.Duration(env.GetEnvAsInt("TIME_SUBTRACTION", 1000)),
		TimeMultiplication: time.Duration(env.GetEnvAsInt("TIME_MULTIPLICATION", 1000)),
		TimeDivision: time.Duration(env.GetEnvAsInt("TIME_DIVISION", 1000)),
	}
}

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

func (a *Orchestrator) RunServer() error {
	r := mux.NewRouter()

	// Добавляем мидлварь для логирования
	r.Use(logger.LoggingMiddleware(a.logger)) 

	a.logger.Info(
		"Оркестратор-Сервер запущен", 
		zap.String("address", fmt.Sprintf(":%s", a.config.OrchestratorPort)),
	)

	return http.ListenAndServe(":"+a.config.OrchestratorPort, r)
}












