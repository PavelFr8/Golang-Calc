package orchestrator

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/PavelFr8/Golang-Calc/pkg/logger"
	"github.com/PavelFr8/Golang-Calc/pkg/tree"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)


type Orchestrator struct {
	Config      *OrchestratorConfig
	logger *zap.Logger
	Expressions   map[string]*Expression
	Tasks   map[string]*Task
	TaskQueue   []*Task
	Mu          sync.Mutex
	exprCounter int64
	taskCounter int64
	
}

func New() *Orchestrator {
	return &Orchestrator{
		Config:    NewOrchestratorConfig(),
		logger: logger.SetupLogger(),
		Expressions: make(map[string]*Expression),
		Tasks: make(map[string]*Task),
		TaskQueue: make([]*Task, 0),
	}
}

func (o *Orchestrator) NewTask(expr *Expression) {
	var traverse func(node *tree.Node)
	traverse = func(node *tree.Node) {
		if node == nil || node.IsLeaf {
			return
		}
		traverse(node.Left)
		traverse(node.Right)
		if node.Left != nil && node.Right != nil && node.Left.IsLeaf && node.Right.IsLeaf {
			if !node.ScheduledTask {
				o.taskCounter++
				taskID := fmt.Sprintf("%d", o.taskCounter)
				var operationTime int
				switch node.Operator {
				case "+":
					operationTime = o.Config.TimeAddition
				case "-":
					operationTime = o.Config.TimeSubtraction
				case "*":
					operationTime = o.Config.TimeMultiplications
				case "/":
					operationTime = o.Config.TimeDivisions
				default:
					operationTime = 100
				}
				task := &Task{
					ID:            taskID,
					ExprID:        expr.ID,
					Arg1:          node.Left.Value,
					Arg2:          node.Right.Value,
					Operation:     node.Operator,
					OperationTime: operationTime,
					Node:          node,
				}
				node.ScheduledTask = true
				o.Tasks[taskID] = task
				o.TaskQueue = append(o.TaskQueue, task)
			}
		}
	}
	traverse(expr.Node)
}

func (o *Orchestrator) RunServer() error {
	r := mux.NewRouter()

	// Добавляем мидлварь для логирования
	r.Use(logger.LoggingMiddleware(o.logger)) 

	r.HandleFunc("/api/v1/calculate", o.CalculateHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/expressions", o.ExpressionsHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/expressions/{id}", o.ExpressionByIDHandler).Methods(http.MethodGet)
	r.HandleFunc("/internal/task", o.GetTaskHandler).Methods(http.MethodGet)
	r.HandleFunc("/internal/task", o.PostTaskHandler).Methods(http.MethodPost)

    r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
    r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("web/js"))))
	r.HandleFunc("/", IndexHandler)

	o.logger.Info(
		"Мой любименький Оркестратор-Сервер, который сжек 50 тысяч моих нервных клеток запущен :)", 
		zap.String("address", fmt.Sprintf(":%s", o.Config.OrchestratorPort)),
	)

	return http.ListenAndServe(":"+o.Config.OrchestratorPort, r)
}
