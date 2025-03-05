package models

import (
	"sync"
	"time"

	"github.com/PavelFr8/Golang-Calc/pkg/env"
)

type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Expression struct {
	ID     int      `json:"id"`
	Status string   `json:"status"`
	Result *float64 `json:"result,omitempty"`
	Tasks  []Task   `json:"-"`
}

var (
	Expressions = make(map[int]*Expression)
	TasksQueue  []Task
	Mutex       sync.Mutex
	ExprCounter int
	TaskCounter int
	TimeAddition =       time.Duration(env.GetEnvAsInt("TIME_ADDITION", 1000))
	TimeSubtraction =   time.Duration(env.GetEnvAsInt("TIME_SUBTRACTION", 1000))
	TimeMultiplication = time.Duration(env.GetEnvAsInt("TIME_MULTIPLICATION", 1000))
	TimeDivision =     time.Duration(env.GetEnvAsInt("TIME_DIVISION", 1000))
)