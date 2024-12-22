package calculation

import "errors"

var (
	ErrInvalidExpression = errors.New("Expression is not valid")
	ErrDivisionByZero    = errors.New("Division by zero")
	ErrCalculation       = errors.New("Error while calculating")
)
