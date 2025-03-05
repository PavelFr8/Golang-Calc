package agent

import (
	"fmt"
	"net/http"

	"github.com/PavelFr8/Golang-Calc/pkg/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

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

	// Запускаем наших демонят(агентов)
	a.StartWorkers()

	return http.ListenAndServe(":"+a.config.AgentPort, r)
}
