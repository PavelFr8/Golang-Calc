package orchestrator

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/PavelFr8/Golang-Calc/pkg/logger"
	"github.com/PavelFr8/Golang-Calc/pkg/tree"
	"github.com/glebarez/sqlite"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
)


type Orchestrator struct {
	Config      *OrchestratorConfig
	logger *zap.Logger
	Expressions   map[uint]*Expression
	Tasks   map[uint]*Task
	TaskQueue   []*Task
	Mu          sync.Mutex
	r *Repository
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&Expression{}, &Task{}); err != nil {
		panic("failed to migrate database")
	}
	return db
}

func New() *Orchestrator {
	return &Orchestrator{
		Config:    NewOrchestratorConfig(),
		logger: logger.SetupLogger(),
		Expressions: make(map[uint]*Expression),
		Tasks: make(map[uint]*Task),
		TaskQueue: make([]*Task, 0),
		r: NewRepository(InitDB()),	
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
				taskID := o.r.GetMaxTaskID() + 1
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
					ExprID:        expr.ID,
					Arg1:          node.Left.Value,
					Arg2:          node.Right.Value,
					Operation:     node.Operator,
					OperationTime: operationTime,
					Node:          node,
				}
				o.r.CreateTask(task)
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
