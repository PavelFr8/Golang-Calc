package agent_test

import (
	"testing"

	"github.com/PavelFr8/Golang-Calc/internal/agent"
	"github.com/PavelFr8/Golang-Calc/internal/agent/models"
	"github.com/stretchr/testify/assert"
)

// Тестирование функции Calc, выполняющей арифметические операции.
func TestCalc(t *testing.T) {
	tests := []struct {
		name     string
		task     models.Task
		expected float64
	}{
		{"addition", models.Task{Arg1: 2, Arg2: 3, Operation: "+"}, 5},
		{"subtraction", models.Task{Arg1: 5, Arg2: 3, Operation: "-"}, 2},
		{"multiplication", models.Task{Arg1: 2, Arg2: 3, Operation: "*"}, 6},
		{"division", models.Task{Arg1: 6, Arg2: 3, Operation: "/"}, 2},
		{"division by zero", models.Task{Arg1: 6, Arg2: 0, Operation: "/"}, 0},
		{"invalid operation", models.Task{Arg1: 2, Arg2: 3, Operation: "^"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.Calc(&tt.task)
			assert.Equal(t, tt.expected, result)
		})	}
}
