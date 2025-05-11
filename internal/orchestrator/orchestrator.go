package orchestrator

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/PavelFr8/Golang-Calc/pkg/env"
	"github.com/PavelFr8/Golang-Calc/pkg/logger"
	"github.com/PavelFr8/Golang-Calc/pkg/tree"
	pb "github.com/PavelFr8/Golang-Calc/proto"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)


type Orchestrator struct {
	Config      *OrchestratorConfig
	logger *zap.Logger
	Expressions   map[uint]*Expression
	Tasks   map[uint]*Task
	TaskQueue   []*Task
	Mu          sync.Mutex
	R *Repository
	pb.OrchestratorServer
}


func New() *Orchestrator {
	return &Orchestrator{
		Config:    NewOrchestratorConfig(),
		logger: logger.SetupLogger(),
		Expressions: make(map[uint]*Expression),
		Tasks: make(map[uint]*Task),
		TaskQueue: make([]*Task, 0),
		R: NewRepository(InitDB()),	
	}
}

func (o *Orchestrator) LoadAndQueuePendingTasks() {
    var exprs []Expression
    err := o.R.DB.Find(&exprs).Error 
	if err != nil {
		if err == gorm.ErrRecordNotFound {
            return
        }
        panic("could not load expressions from DB: " + err.Error())
    }

    for i := range exprs {
        e := &exprs[i]
        node, err := tree.BuildNode(e.Expr)
        if err != nil {
            panic("FAIL TO BUILD NODE" + err.Error())
        }
        e.Node = node
        o.Expressions[e.ID] = e
		if e.Status == "pending" {
			o.NewTask(e)
		}
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
				taskID := o.R.GetMaxTaskID() + 1
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
				o.R.CreateTask(task)
				node.ScheduledTask = true
				o.Tasks[taskID] = task
				o.TaskQueue = append(o.TaskQueue, task)
			}
		}
	}
	traverse(expr.Node)
}

func (o *Orchestrator) RunGrpc() {
	addr := env.GetEnv("ORCHESTRATOR_GRPC_ADDR", "localhost:5000")
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		o.logger.Info("error starting tcp listener: " + err.Error())
		os.Exit(1)
	}
	
	o.logger.Info("tcp listener started at address: " + addr)
	grpcServer := grpc.NewServer()
	pb.RegisterOrchestratorServer(grpcServer, o)
	// запустим grpc сервер
	if err := grpcServer.Serve(lis); err != nil {
		o.logger.Info("error serving grpc: " + err.Error())
		os.Exit(1)
	}
}

func (o *Orchestrator) RunServer() error {
	r := mux.NewRouter()

	// Добавляем мидлварь для логирования
	r.Use(logger.LoggingMiddleware(o.logger))

	exempt := map[string]bool{
		"/api/v1/login":    true,
		"/api/v1/register": true,
		"/":                true,
	}

	r.Use(JWTMiddleware(o.Config.JWTsecret, exempt))

	r.HandleFunc("/api/v1/calculate", o.CalculateHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/expressions", o.ExpressionsHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/expressions/{id}", o.ExpressionByIDHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/register", o.RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/login", o.LoginHandler).Methods(http.MethodPost)

    r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
    r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("web/js"))))
	r.HandleFunc("/", IndexHandler)

	o.LoadAndQueuePendingTasks()

	o.logger.Info(
		"Мой любименький Оркестратор-Сервер, который сжек 50 тысяч моих нервных клеток запущен :)", 
		zap.String("address", fmt.Sprintf(":%s", o.Config.OrchestratorPort)),
	)

	go o.RunGrpc()

	return http.ListenAndServe(":"+o.Config.OrchestratorPort, r)
}
