package orchestrator

import (
	"github.com/PavelFr8/Golang-Calc/pkg/tree"
)

type Task struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	ExprID        uint         `json:"-"`
	Arg1          *float64     `json:"arg1"`
	Arg2          *float64     `json:"arg2"`
	Operation     string       `json:"operation"`
	OperationTime int          `json:"operation_time"`
	Result        *float64     `json:"-" gorm:"default:(-)"`
	Node          *tree.Node   `gorm:"-" json:"-"`
}

type Expression struct {
	ID       uint        `gorm:"primaryKey" json:"id"`
	Expr     string      `json:"expression"`
	Status   string      `json:"status"`
	Result   *float64    `json:"result,omitempty"`
	Tasks    []Task      `gorm:"foreignKey:ExprID" json:"tasks,omitempty"`
	Node     *tree.Node  `gorm:"-" json:"-"`
}