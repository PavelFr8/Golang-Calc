package tree_test

import (
	"testing"

	"github.com/PavelFr8/Golang-Calc/pkg/tree"
)

func TestBuildNode_SimpleExpression(t *testing.T) {
	node, err := tree.BuildNode("2 + 3")
	if err != nil {
		t.Fatalf("Ожидалось успешное построение дерева, получена ошибка: %v", err)
	}
	if node == nil || node.Operator != "+" {
		t.Errorf("Ожидался оператор '+', получено: %+v", node)
	}
	if node.Left == nil || *node.Left.Value != 2 {
		t.Errorf("Ожидалось левое значение 2, получено: %+v", node.Left)
	}
	if node.Right == nil || *node.Right.Value != 3 {
		t.Errorf("Ожидалось правое значение 3, получено: %+v", node.Right)
	}
}

func TestBuildNode_InvalidExpression(t *testing.T) {
	_, err := tree.BuildNode("2 +")
	if err == nil {
		t.Error("Ожидалась ошибка при некорректном выражении, но её не произошло")
	}
}

func TestBuildNode_ComplexExpression(t *testing.T) {
	node, err := tree.BuildNode("(1 + 2) * 3")
	if err != nil {
		t.Fatalf("Ожидалось успешное построение дерева, ошибка: %v", err)
	}
	if node.Operator != "*" {
		t.Errorf("Ожидался оператор '*', получено: %s", node.Operator)
	}
	if node.Left == nil || node.Left.Operator != "+" {
		t.Errorf("Левое поддерево должно быть сложением, получено: %+v", node.Left)
	}
	if node.Right == nil || *node.Right.Value != 3 {
		t.Errorf("Правое значение должно быть 3, получено: %+v", node.Right)
	}
}
