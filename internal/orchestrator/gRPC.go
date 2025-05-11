package orchestrator

import (
	"context"
	"fmt"

	pb "github.com/PavelFr8/Golang-Calc/proto"
)

// Тупо отдаем самый последний элемент очереди
func (o *Orchestrator) GetTask(ctx context.Context, in *pb.Empty) (*pb.Task, error) {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if len(o.TaskQueue) == 0 {
		return nil, fmt.Errorf("No task available")
	}

	task := o.TaskQueue[0]
	o.TaskQueue = o.TaskQueue[1:]

	if _, exists := o.Expressions[task.ExprID]; !exists {
		return nil, fmt.Errorf("Task expression not found")
	}
	grpc_task := &pb.Task{
		ID: uint32(task.ID),
		Arg1: *task.Arg1,
		Arg2: *task.Arg2,
		Operation: task.Operation,
		OperationTime: int32(task.OperationTime),
	}
	return grpc_task, nil
}

// Тут уже с огромной болью с слезами, добавляем решенное выражение обратно в очередь
func (o *Orchestrator) PostTask(ctx context.Context, grpc_task *pb.TaskResult) (*pb.Empty, error) {
	o.Mu.Lock()
	task, ok := o.Tasks[uint(grpc_task.ID)]
	if !ok {
		o.Mu.Unlock()
		return nil, fmt.Errorf("Task not found")
	}
	task.Result = &grpc_task.Result
	task.Node.IsLeaf = true
	task.Node.Value = &grpc_task.Result
	delete(o.Tasks, uint(grpc_task.ID))
	o.r.db.Updates(task)
	if expr, exists := o.Expressions[task.ExprID]; exists {
		o.NewTask(expr)
		if expr.Node.IsLeaf {
			expr.Status = "completed"
			expr.Result = expr.Node.Value
			o.r.db.Updates(expr)
		}
	}
	o.Mu.Unlock()
	return nil, nil
}
