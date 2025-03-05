package agent

import (
	"github.com/PavelFr8/Golang-Calc/pkg/env"
)

type AgentConfig struct {
	AgentPort string
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{
		AgentPort: env.GetEnv("AGENT_PORT", "8081"),
	}
}