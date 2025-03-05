package tree

import "errors"

var (
	ErrConsecutiveOperators       = errors.New("бРАТУХА, два оператора подряд")
	ErrDivisionByZero             = errors.New("Емае, на ноль делить нельзя")
	ErrEmptyExpression            = errors.New("А где? А где выражение?")
	ErrFormedExpression           = errors.New("Что за выражение такое?")
	ErrUnsupportedOperator        = errors.New("Самый умный, да? Не поддерживаемый оператор")
	ErrInvalidCharacterInInput    = errors.New("НУ ТЫ БАЛБЕС! Недопустимый символ в выражении")
	ErrExpressionEndsWithOperator = errors.New("Математик фигов, выражение заканчивается оператором у тебя!")
)
