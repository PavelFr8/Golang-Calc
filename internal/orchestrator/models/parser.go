package models

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"time"
)

func ParseExpression(exprID int, exprStr string) (*Expression, error) {
	exprAST, err := parser.ParseExpr(exprStr)
	if err != nil {
		return nil, err
	}

	expression := &Expression{
		ID:     exprID,
		Status: "cooking...(in progress, wait, bro)",
	}

	ast.Inspect(exprAST, func(n ast.Node) bool {
		if binExpr, ok := n.(*ast.BinaryExpr); ok {
			Mutex.Lock()
			TaskCounter++
			taskID := TaskCounter
			Mutex.Unlock()

			task := Task{
				ID:            taskID,
				Arg1:          extractValue(binExpr.X),
				Arg2:          extractValue(binExpr.Y),
				Operation:     binExpr.Op.String(),
				OperationTime: int(getOperationTime(binExpr.Op)),
			}
			expression.Tasks = append(expression.Tasks, task)
		}
		return true
	})

	return expression, nil
}

func extractValue(node ast.Expr) float64 {
	if lit, ok := node.(*ast.BasicLit); ok {
		val, _ := strconv.ParseFloat(lit.Value, 64)
		return val
	}
	return 0
}

func getOperationTime(op token.Token) time.Duration {
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
		return 1000 * time.Millisecond
	}
}