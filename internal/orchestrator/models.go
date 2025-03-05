package orchestrator

import (
	"sync"

	"github.com/PavelFr8/Golang-Calc/pkg/env"
	"github.com/PavelFr8/Golang-Calc/pkg/tree"
)

type Task struct {
	ID            string       `json:"id"`
	ExprID        string       `json:"-"`
	Arg1          float64      `json:"arg1"`
	Arg2          float64      `json:"arg2"`
	Operation     string       `json:"operation"`
	OperationTime int          `json:"operation_time"`
	Node          *tree.Node `json:"-"`
}

type Expression struct {
	ID     string       `json:"id"`
	Expr   string       `json:"expression"`
	Status string       `json:"status"`
	Result *float64     `json:"result,omitempty"`
	Node    *tree.Node `json:"-"`
}

var (
	Expressions = make(map[int]*Expression)
	TasksQueue  []Task
	Mutex       sync.Mutex
	ExprCounter int
	TaskCounter int
	TimeAddition =       env.GetEnvAsInt("TIME_ADDITION_MS", 1000)
	TimeSubtraction =   env.GetEnvAsInt("TIME_SUBTRACTION_MS", 1000)
	TimeMultiplication = env.GetEnvAsInt("TIME_MULTIPLICATIONS_MS", 1000)
	TimeDivision =     env.GetEnvAsInt("TIME_DIVISIONS_MS", 1000)
)