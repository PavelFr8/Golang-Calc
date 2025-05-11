package orchestrator

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/PavelFr8/Golang-Calc/pkg/logger"
	"github.com/PavelFr8/Golang-Calc/pkg/tree"
	pb "github.com/PavelFr8/Golang-Calc/proto"
	"github.com/glebarez/sqlite"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	pb.OrchestratorServer
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

func (o *Orchestrator) LoadAndQueuePendingTasks() {
    var exprs []Expression
    err := o.r.db.Find(&exprs).Error 
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

func (o *Orchestrator) RunGrpc() {
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		o.logger.Info("error starting tcp listener: " + err.Error())
		os.Exit(1)
	}
	
	o.logger.Info("tcp listener started at port: " + port)
	grpcServer := grpc.NewServer()
	// зарегистрируем нашу реализацию сервера
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

	r.HandleFunc("/api/v1/calculate", o.CalculateHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/expressions", o.ExpressionsHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/expressions/{id}", o.ExpressionByIDHandler).Methods(http.MethodGet)

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
