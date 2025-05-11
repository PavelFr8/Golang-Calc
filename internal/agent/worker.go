package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/PavelFr8/Golang-Calc/pkg/env"
	pb "github.com/PavelFr8/Golang-Calc/proto"
	"go.uber.org/zap"
)

var (
	computingPower  = env.GetEnvAsInt("COMPUTING_POWER", 4)
	demonSleepTime = time.Duration(5000 * time.Millisecond)
	wg sync.WaitGroup
)

func Calc(task *pb.Task) float64 {
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		return task.Arg1 / task.Arg2
	default:
		return 0
	}
}

// Демон, который получает выражение для вычисления с оркестратора, вычисляет его и отправляет на оркестратор результат выражения.
func (a *Agent) Worker() {
	defer wg.Done()
	for {
		task, err := a.grpcClient.GetTask(context.TODO(), nil)
		if err != nil {
			fmt.Println("Ошибка получения задачи:", err)
			time.Sleep(demonSleepTime)
			continue
		}
		if task == nil {
			time.Sleep(demonSleepTime)
			continue
		}

		fmt.Printf("Получена задача: %v\n", *task)
		result := Calc(task)
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		resultTask := &pb.TaskResult{
			ID: task.ID,
			Result: result,
		}

		if _, err := a.grpcClient.PostTask(context.TODO(), resultTask); err != nil {
			fmt.Println("Ошибка отправки результата:", err)
		} else {
			fmt.Printf("Задача %d выполнена, результат: %f\n", task.ID, result)
		}
	}
}

func (a *Agent) StartWorkers() {
	a.logger.Info("Запуск супер-секретный-демонов-агентов", zap.Int("computingPower", computingPower))
	for i := 0; i < computingPower; i++ {
		wg.Add(1)
		go a.Worker()
	}
	wg.Wait()
}
