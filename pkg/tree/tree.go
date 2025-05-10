package tree

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Короче, пакет для работы с деревьями AST, который используется в калькуляторе.

type Node struct {
	IsLeaf        bool
	Value         float64
	Operator      string
	Left, Right   *Node
	ScheduledTask bool
}

// функция, которая строит дерево AST из строки с математическим выражением.
func BuildNode(expr string) (*Node, error) {
	expr = strings.ReplaceAll(expr, " ", "")
	if expr == "" {
		return nil, ErrEmptyExpression
	}
	p := &parser{input: expr, index: 0}
	ast, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if p.index < len(p.input) {
		return nil, fmt.Errorf("unexpected token at position %d", p.index)
	}
	return ast, nil
}

type parser struct {
	input string
	index int
}

// функция, которая сдвигает указатель на один символ вперед.
func (p *parser) advance() rune {
	ch := p.peek()
	p.index++
	return ch
}

// функция, которая возвращает текущий символ в строке.
func (p *parser) peek() rune {
	if p.index < len(p.input) {
		return rune(p.input[p.index])
	}
	return 0
}

// функция, которая парсит фактор.
func (p *parser) parseElement() (*Node, error) {
	ch := p.peek()
	if ch == '(' {
		p.advance()
		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.peek() != ')' {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		p.advance()
		return node, nil
	}
	start := p.index
	if ch == '+' || ch == '-' {
		p.advance()
	}
	for {
		ch = p.peek()
		if unicode.IsDigit(ch) || ch == '.' {
			p.advance()
		} else {
			break
		}
	}
	token := p.input[start:p.index]
	if token == "" {
		return nil, fmt.Errorf("expected number at position %d", start)
	}
	value, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number %s", token)
	}
	return &Node{
		IsLeaf: true,
		Value:  value,
	}, nil
}

// функция, которая тупо парсит выражение :)
func (p *parser) parseExpression() (*Node, error) {
	node, err := p.parseSubExpr()
	if err != nil {
		return nil, err
	}
	for {
		ch := p.peek()
		if ch == '+' || ch == '-' {
			op := string(p.advance())
			rightNode, err := p.parseSubExpr()
			if err != nil {
				return nil, err
			}
			node = &Node{
				IsLeaf:   false,
				Operator: op,
				Left:     node,
				Right:    rightNode,
			}
		} else {
			break
		}
	}
	return node, nil
}

// функция, которая парсит терм.
func (p *parser) parseSubExpr() (*Node, error) {
	node, err := p.parseElement()
	if err != nil {
		return nil, err
	}
	for {
		ch := p.peek()
		if ch == '*' || ch == '/' {
			op := string(p.advance())
			rightNode, err := p.parseElement()
			if rightNode.Value == 0 && ch == '/' {
				return nil, ErrDivisionByZero
			}
			if err != nil {
				return nil, err
			}
			node = &Node{
				IsLeaf:   false,
				Operator: op,
				Left:     node,
				Right:    rightNode,
			}
		} else {
			break
		}
	}
	return node, nil
}
