package main

import (
	"github.com/PavelFr8/Golang-Calc/internal/orchestrator"
)

func main() {
	orchestrator := orchestrator.New()
	orchestrator.RunServer()
}
