// test/agent/agent_mock_test.go
package agent_test

import (
	"context"
	"net"
	"testing"
	"time"

	agentpkg "github.com/PavelFr8/Golang-Calc/internal/agent"
	pb "github.com/PavelFr8/Golang-Calc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type mockOrchestratorServer struct {
	pb.UnimplementedOrchestratorServer
	t         *testing.T
	taskSent  bool
	taskDone  bool
}

func (m *mockOrchestratorServer) GetTask(ctx context.Context, _ *pb.Empty) (*pb.Task, error) {
	if m.taskSent {
		return nil, nil // имитируем отсутствие задач
	}
	m.taskSent = true
	return &pb.Task{
		ID:            1,
		Arg1:          4,
		Arg2:          5,
		Operation:     "*",
		OperationTime: 10,
	}, nil
}

func (m *mockOrchestratorServer) PostTask(ctx context.Context, res *pb.TaskResult) (*pb.Empty, error) {
	m.taskDone = true
	if res.ID != 1 || res.Result != 20 {
		m.t.Errorf("Ожидалось ID=1 и Result=20, получено ID=%d Result=%f", res.ID, res.Result)
	} else {
		m.t.Logf("✅ Агент корректно отправил результат: %v", res.Result)
	}
	return &pb.Empty{}, nil
}

func TestAgentWithMockOrchestrator(t *testing.T) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		t.Fatalf("не удалось запустить сервер: %v", err)
	}

	grpcServer := grpc.NewServer()
	mockServer := &mockOrchestratorServer{t: t}
	pb.RegisterOrchestratorServer(grpcServer, mockServer)
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	time.Sleep(200 * time.Millisecond)

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("ошибка подключения агента: %v", err)
	}
	client := pb.NewOrchestratorClient(conn)
	a := &agentpkg.Agent{
		GrpcClient: client,
	}

	// однократное выполнение воркера
	oneShotWorker(t, a)

	if !mockServer.taskDone {
		t.Fatal("Агент не завершил обработку задачи")
	}
}

func oneShotWorker(t *testing.T, a *agentpkg.Agent) {
	task, err := a.GrpcClient.GetTask(context.TODO(), &pb.Empty{})
	if err != nil || task == nil {
		t.Fatalf("не удалось получить задачу: %v", err)
	}
	result := agentpkg.Calc(task)
	res := &pb.TaskResult{ID: task.ID, Result: result}
	_, err = a.GrpcClient.PostTask(context.TODO(), res)
	if err != nil {
		t.Fatalf("не удалось отправить результат: %v", err)
	}
}