package orchestrator

import (
	"github.com/PavelFr8/Golang-Calc/pkg/env"
)

type OrchestratorConfig struct {
	OrchestratorPort    string
	JWTsecret           []byte
	TimeAddition        int
	TimeSubtraction     int
	TimeMultiplications int
	TimeDivisions       int
}

func NewOrchestratorConfig() *OrchestratorConfig {
	return &OrchestratorConfig{
		OrchestratorPort:    env.GetEnv("ORCHESTRATOR_PORT", "8080"),
		JWTsecret:           []byte(env.GetEnv("JWT_SECRET", "qwerty12345")),
		TimeAddition:        env.GetEnvAsInt("TIME_ADDITION_MS", 1000),
		TimeSubtraction:     env.GetEnvAsInt("TIME_SUBTRACTION_MS", 1000),
		TimeMultiplications: env.GetEnvAsInt("TIME_MULTIPLICATIONS_MS", 1000),
		TimeDivisions:       env.GetEnvAsInt("TIME_DIVISIONS_MS", 1000),}
}