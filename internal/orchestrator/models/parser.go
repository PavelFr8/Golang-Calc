package models

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func parseAST(expression *Expression, node ast.Expr) float64 {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		// Рекурсивно вычисляем операнды
		left := parseAST(expression, n.X)
		right := parseAST(expression, n.Y)

		// Генерируем уникальный ID для задачи
		Mutex.Lock()
		TaskCounter++
		taskID := TaskCounter
		

		// Создание задачи
		task := Task{
			ID:            taskID,
			Arg1:          left,
			Arg2:          right,
			Operation:     n.Op.String(),
			OperationTime: getOperationTime(n.Op),
		}
		expression.Tasks = append(expression.Tasks, task)
		Mutex.Unlock()

		// Возвращаем ID задачи для обработки агентом
		return float64(taskID)

	case *ast.BasicLit:
		// Листовые узлы, содержащие числа
		val, _ := strconv.ParseFloat(n.Value, 64)
		return val

	case *ast.ParenExpr: // Обработка скобок
		return parseAST(expression, n.X)

	default:
		return 0
	}
}

// Функция для парсинга выражения
func ParseExpression(exprID int, exprStr string) (*Expression, error) {
	exprAST, err := parser.ParseExpr(exprStr)
	if err != nil {
		return nil, err
	}

	expression := &Expression{
		ID:     exprID,
		Status: "cooking...(in progress, wait, bro)",
	}

	// Парсим AST и генерируем задачи
	parseAST(expression, exprAST)

	return expression, nil
}

// Получаем время для операции, указанное в переменных среды
func getOperationTime(op token.Token) int {
	switch op {
	case token.ADD:
		return TimeAddition
	case token.SUB:
		return TimeSubtraction
	case token.MUL:
		return TimeMultiplication
	case token.QUO:
		return TimeDivision
	default:
		return 1000
	}
}