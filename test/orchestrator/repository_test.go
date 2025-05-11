// test/orchestrator/repository_test.go
package orchestrator_test

import (
	"testing"

	"github.com/PavelFr8/Golang-Calc/internal/orchestrator"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestRepo(t *testing.T) *orchestrator.Repository {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("не удалось подключиться к тестовой БД: %v", err)
	}
	if err := db.AutoMigrate(&orchestrator.Expression{}, &orchestrator.Task{}, &orchestrator.User{}); err != nil {
		t.Fatalf("не удалось выполнить миграции: %v", err)
	}
	return orchestrator.NewRepository(db)
}

func TestCreateUser(t *testing.T) {
	repo := setupTestRepo(t)
	err := repo.CreateUser("testuser", "password")
	if err != nil {
		t.Fatalf("ошибка при создании пользователя: %v", err)
	}

	var user orchestrator.User
	err = repo.DB.First(&user, "login = ?", "testuser").Error
	if err != nil {
		t.Fatalf("пользователь не найден в БД: %v", err)
	}
	if user.Login != "testuser" {
		t.Errorf("ожидался логин 'testuser', получено: %s", user.Login)
	}
}

func TestCreateExpressionAndTask(t *testing.T) {
	repo := setupTestRepo(t)
	expr := &orchestrator.Expression{
		Expr:   "1 + 2",
		Status: "pending",
	}
	err := repo.CreateExpression(expr)
	if err != nil {
		t.Fatalf("не удалось создать выражение: %v", err)
	}

	task := &orchestrator.Task{
		ExprID:        expr.ID,
		Arg1:          float64Ptr(1),
		Arg2:          float64Ptr(2),
		Operation:     "+",
		OperationTime: 100,
	}
	err = repo.CreateTask(task)
	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	var savedTask orchestrator.Task
	err = repo.DB.First(&savedTask, task.ID).Error
	if err != nil {
		t.Fatalf("задача не найдена: %v", err)
	}
	if savedTask.Operation != "+" || *savedTask.Arg1 != 1 || *savedTask.Arg2 != 2 {
		t.Errorf("неверные данные задачи: %+v", savedTask)
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestGetMaxIDs(t *testing.T) {
	repo := setupTestRepo(t)
	if id := repo.GetMaxExpressionID(); id != 1 {
		t.Errorf("ожидался id=1, получено: %d", id)
	}
	repo.CreateExpression(&orchestrator.Expression{Expr: "a"})
	repo.CreateExpression(&orchestrator.Expression{Expr: "b"})
	if id := repo.GetMaxExpressionID(); id != 3 {
		t.Errorf("ожидался id=3, получено: %d", id)
	}
}