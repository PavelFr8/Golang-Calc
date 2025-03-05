package agent

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PavelFr8/Golang-Calc/pkg/env"
	"github.com/PavelFr8/Golang-Calc/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AgentConfig struct {
	AgentPort string
	OrchestratorUrl string
	TimeAddition time.Duration
	TimeSubtraction time.Duration
	TimeMultiplication time.Duration
	TimeDivision time.Duration
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		AgentPort: env.GetEnv("AGENT_PORT", "8081"),
		OrchestratorUrl: env.GetEnv("ORCHESTRATOR_URL", "8080"),
		TimeAddition: time.Duration(env.GetEnvAsInt("TIME_ADDITION", 1000)),
		TimeSubtraction: time.Duration(env.GetEnvAsInt("TIME_SUBTRACTION", 1000)),
		TimeMultiplication: time.Duration(env.GetEnvAsInt("TIME_MULTIPLICATION", 1000)),
		TimeDivision: time.Duration(env.GetEnvAsInt("TIME_DIVISION", 1000)),
	}
}

type Agent struct {
	config *AgentConfig
	logger *zap.Logger
}

func New() *Agent {
	return &Agent{
		config: NewAgentConfig(),
		logger: logger.SetupLogger(),
	}
}

func (a *Agent) RunServer() error {
	r := mux.NewRouter()

	// Добавляем мидлварь для логирования
	r.Use(logger.LoggingMiddleware(a.logger)) 

	a.logger.Info(
		"Агент-Сервер запущен", 
		zap.String("address", fmt.Sprintf(":%s", a.config.AgentPort)),
	)

	return http.ListenAndServe(":"+a.config.AgentPort, r)
}
