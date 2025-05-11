package agent_test

import (
	"testing"

	"github.com/PavelFr8/Golang-Calc/internal/agent"
	pb "github.com/PavelFr8/Golang-Calc/proto"
	"github.com/stretchr/testify/assert"
)

// Тестирование функции Calc, выполняющей арифметические операции.
func TestCalc(t *testing.T) {
	tests := []struct {
		name     string
		task     pb.Task
		expected float64
	}{
		{"addition", pb.Task{Arg1: 2, Arg2: 3, Operation: "+"}, 5},
		{"subtraction", pb.Task{Arg1: 5, Arg2: 3, Operation: "-"}, 2},
		{"multiplication", pb.Task{Arg1: 2, Arg2: 3, Operation: "*"}, 6},
		{"division", pb.Task{Arg1: 6, Arg2: 3, Operation: "/"}, 2},
		{"division by zero", pb.Task{Arg1: 6, Arg2: 0, Operation: "/"}, 0},
		{"invalid operation", pb.Task{Arg1: 2, Arg2: 3, Operation: "^"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.Calc(&tt.task)
			assert.Equal(t, tt.expected, result)
		})	}
}
