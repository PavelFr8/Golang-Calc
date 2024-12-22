package main

import (
	"github.com/PavelFr8/Golang-Calc/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
