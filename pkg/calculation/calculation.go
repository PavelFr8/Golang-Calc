package calculation

import (
	"slices"
	"strconv"
	"strings"
)

// Проверка выражения на корректность написания
func is_valid_expression(expression string) bool {
	if !(strings.Count(expression, "(") == strings.Count(expression, ")")) {
		return false // Ошибка в скобочках
	}
	correct_symbols := []string{"+", "-", "*", "/", "(", ")", "."}
	for _, r := range expression {
		if !strings.Contains("0123456789", string(r)) && !slices.Contains(correct_symbols, string(r)) {
			return false // Использованы некорректные символы
		}
	}
	return true
}

// Конвертер в обратную польскую последовательность
func convert_to_polish(str string) string {
	priority := map[string]int{
		"+": 1,
		"-": 1,
		"/": 2,
		"*": 2,
	}

	var result []string
	var stack []rune
	var number strings.Builder

	for _, r := range str {
		if strings.Contains("0123456789.", string(r)) {
			number.WriteRune(r) // Добавляем символ к числу
		} else {
			if number.Len() > 0 {
				result = append(result, number.String()) // Если число было собрано, добавляем его в результат
				number.Reset()                           // Сбрасываем сборщик для следующего числа
			}

			if r == '(' {
				stack = append(stack, r)
			} else if r == ')' {
				for len(stack) > 0 && stack[len(stack)-1] != '(' {
					result = append(result, string(stack[len(stack)-1]))
					stack = stack[:len(stack)-1]
				}
				if len(stack) > 0 {
					stack = stack[:len(stack)-1] // Убираем '('
				}
			} else {
				for len(stack) > 0 && priority[string(stack[len(stack)-1])] >= priority[string(r)] {
					result = append(result, string(stack[len(stack)-1]))
					stack = stack[:len(stack)-1]
				}
				stack = append(stack, r)
			}
		}
	}

	// Добавляем последнее число, если оно есть
	if number.Len() > 0 {
		result = append(result, number.String())
	}

	for len(stack) > 0 {
		result = append(result, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return strings.Join(result, " ")
}

// Простая функция для вычисления
func do_calculation(n1 float64, n2 float64, operation string) (float64, error) {
	switch operation {
	case "+":
		return n1 + n2, nil
	case "-":
		return n1 - n2, nil
	case "*":
		return n1 * n2, nil
	case "/":
		if n2 == 0 {
			return 0, ErrDivisionByZero
		}
		return n1 / n2, nil
	default:
		return 0, ErrCalculation
	}
}

// Парсим строку в формате ОПС и проводим вычисления
func calculate(str string) (float64, error) {
	var stack []float64

	for _, elem := range strings.Fields(str) {
		if num, err := strconv.ParseFloat(elem, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}

			n1 := stack[0]
			n2 := stack[1]
			stack = stack[:len(stack)-2]

			result, err := do_calculation(n1, n2, elem)
			if err != nil {
				return 0, err
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, ErrCalculation
	}
	return stack[0], nil
}

// Основная функция, грубо говоря - калькулятор
func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "") // Удаляем пробелы

	// Проверка выражения на корректность
	if !is_valid_expression(expression) {
		return 0, ErrInvalidExpression
	}

	polish := convert_to_polish(expression)
	result, err := calculate(polish)

	return result, err
}
