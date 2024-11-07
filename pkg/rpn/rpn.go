package rpn

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

// проверка на корректное написание скобок
func check_scobs(str string) bool {
	return strings.Count(str, "(") == strings.Count(str, ")")
}

// проверка на то, что во входных данных только разрешенные символы
func check_elems(str string) bool {
	good_boys := []string{"+", "-", "*", "/", "(", ")", "."}
	for _, r := range str {
		if !strings.Contains("0123456789", string(r)) && !slices.Contains(good_boys, string(r)) {
			return false
		}
	}
	return true
}

// конвертер в обратную польскую последовательность
func convert_to_polish(str string) string {
	priority := map[string]int{
		"+": 1,
		"-": 1,
		"/": 2,
		"*": 2,
	}

	var result []string
	var stack []rune

	for _, r := range str {
		if strings.Contains("0987654321", string(r)) {
			result = append(result, string(r))
		} else if r == '(' {
			stack = append(stack, r)
		} else if r == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				result = append(result, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		} else {
			for len(stack) > 0 && priority[string(stack[len(stack)-1])] >= priority[string(r)] {
				result = append(result, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, r)
		}
	}

	for len(stack) > 0 {
		result = append(result, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return strings.Join(result, " ")
}

// простая функция для вычислений
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
			return 0, errors.New("На ноль делить нельзя!")
		}
		return n1 / n2, nil
	default:
		return 0, errors.New("Ошибочка в вычислениях?")
	}
}

// парсим строку в формате ОПС и проводим вычисления
func calculate(str string) (float64, error) {
	var stack []float64

	for _, elem := range strings.Fields(str) {
		if num, err := strconv.ParseFloat(elem, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("Ошибка при конвертации в обратную польскую последовательность. Ошибка входных данных")
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
		return 0, errors.New("Ошибка при выполнении вычислений")
	}
	return stack[0], nil
}

// основная функция, грубо говоря - калькулятор
func Calc(expression string) (float64, error) {
	if !check_scobs(expression) {
		return 0, errors.New("Чё по скобочкам?!")
	}

	expression = strings.ReplaceAll(expression, " ", "")

	if !check_elems(expression) {
		return 0, errors.New("Чё по входным данным!?")
	}

	polish := convert_to_polish(expression)

	anw, err := calculate(polish)

	return anw, err
}
