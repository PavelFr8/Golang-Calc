package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/PavelFr8/Golang-Calc/internal/agent/models"
	"github.com/PavelFr8/Golang-Calc/pkg/env"
	"go.uber.org/zap"
)

var (
	orchestratorURL = env.GetEnv("ORCHESTRATOR_URL", "http://localhost:8081")
	computingPower  = env.GetEnvAsInt("COMPUTING_POWER", 4)
	demonSleepTime = time.Duration(5000 * time.Millisecond)
	wg sync.WaitGroup
)

func getTask() (*models.Task, error) {
	resp, err := http.Get(orchestratorURL + "/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	var task struct {
		Task models.Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task.Task, nil
}


func calc(task *models.Task) float64 {
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 != 0 {
			return task.Arg1 / task.Arg2
		}
		return 0
	default:
		return 0
	}
}

// Функция для отправки выражения на оркестратор
func submitResult(taskID int, result float64) error {
	resultPayload := models.TaskResult{
		ID:     taskID,
		Result: result,
	}
	data, err := json.Marshal(resultPayload)
	if err != nil {
		return err
	}

	resp, err := http.Post(orchestratorURL+"/internal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to submit result, status code: %d", resp.StatusCode)
	}
	return nil
}

// Демон, который получает выражение для вычисления с оркестратора, вычисляет его и отправляет на оркестратор результат выражения.
func worker() {
	defer wg.Done()
	for {
		task, err := getTask()
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
		result := calc(task)
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		fmt.Println(time.Duration(task.OperationTime) * time.Millisecond)

		if err := submitResult(task.ID, result); err != nil {
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
		go worker()
	}
	wg.Wait()
}
