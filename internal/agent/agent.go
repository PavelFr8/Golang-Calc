package agent

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PavelFr8/Golang-Calc/pkg/env"
	"github.com/PavelFr8/Golang-Calc/pkg/logger"
	pb "github.com/PavelFr8/Golang-Calc/proto"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Agent struct {
	config *AgentConfig
	logger *zap.Logger
	GrpcClient pb.OrchestratorClient
}

func New() *Agent {
	return &Agent{
		config: NewAgentConfig(),
		logger: logger.SetupLogger(),
		GrpcClient: ConnectToGrpcService(),
	}
}

func ConnectToGrpcService() pb.OrchestratorClient {
	addr := env.GetEnv("ORCHESTRATOR_GRPC_ADDR", "localhost:5000")
	fmt.Println("Connecting to gRPC server at:", addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("could not connect to grpc server: " + err.Error())
		os.Exit(1)
	}
	// defer conn.Close()
	grpcClient := pb.NewOrchestratorClient(conn)
	return grpcClient
}

func (a *Agent) RunServer() error {
	r := mux.NewRouter()

	// Добавляем мидлварь для логирования
	r.Use(logger.LoggingMiddleware(a.logger)) 

	a.logger.Info(
		"Агент-Сервер запущен", 
		zap.String("address", fmt.Sprintf(":%s", a.config.AgentPort)),
	)

	// Запускаем наших демонят(агентов)
	a.StartWorkers()

	return http.ListenAndServe(":"+a.config.AgentPort, r)
}
