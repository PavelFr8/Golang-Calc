package main

import (
	"github.com/PavelFr8/Golang-Calc/internal/agent"
)

func main() {
	agent := agent.New()
	agent.RunServer()
}
