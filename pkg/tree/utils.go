package tree

import "fmt"

func Check(node *Node) error {
	if node == nil {
		return ErrEmptyExpression
	}
	if err := FullCheck(node); err != nil {
		return err
	}
	return nil
}

func FullCheck(node *Node) error {
	if node.IsLeaf {
		return nil
	}

	if node.Operator != "+" && node.Operator != "-" && node.Operator != "/" && node.Operator != "*" {
		return ErrUnsupportedOperator
	}

	if node.Left == nil || node.Right == nil {
		return fmt.Errorf("missing operand for operator %s", node.Operator)
	}

	if node.Left != nil && !node.Left.IsLeaf {
		if node.Left.Operator != "" {
			return ErrConsecutiveOperators
		}
	}
	if node.Right != nil && !node.Right.IsLeaf {
		if node.Right.Operator != "" {
			return ErrConsecutiveOperators
		}
	}

	if node.Operator == "/" && node.Right != nil && node.Right.Value == 0 {
		return ErrDivisionByZero
	}

	if err := FullCheck(node.Left); err != nil {
		return err
	}
	if err := FullCheck(node.Right); err != nil {
		return err
	}

	return nil
}