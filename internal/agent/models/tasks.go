package models

type Task struct {
	ID            uint    `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type TaskResult struct {
	ID     uint    `json:"id"`
	Result float64 `json:"result"`
}