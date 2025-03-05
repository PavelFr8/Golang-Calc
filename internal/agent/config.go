package agent

import (
	"github.com/PavelFr8/Golang-Calc/pkg/env"
)

type AgentConfig struct {
	AgentPort string
	OrchestratorURL string
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		AgentPort: env.GetEnv("AGENT_PORT", "8081"),
		OrchestratorURL: env.GetEnv("ORCHESTRATOR_URL", "http://localhost:8080"),
	}
}