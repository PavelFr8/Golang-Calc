package orchestrator

import (
	"github.com/PavelFr8/Golang-Calc/pkg/env"
)

type OrchestratorConfig struct {
	OrchestratorPort   string
}

func NewOrchestratorConfig() *OrchestratorConfig {
	return &OrchestratorConfig{
		OrchestratorPort:   env.GetEnv("ORCHESTRATOR_PORT", "8081"),
	}
}
