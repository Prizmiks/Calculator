package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var precedence = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}

func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	if len(tokens) == 0 {
		return 0, errors.New("пустое выражение")
	}

	postfix, err := infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}

	return evaluatePostfix(postfix)
}

func tokenize(expression string) []string {
	var tokens []string
	var number strings.Builder

	for _, char := range expression {
		if char >= '0' && char <= '9' || char == '.' {
			number.WriteRune(char)
		} else {
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			if char != ' ' {
				tokens = append(tokens, string(char))
			}
		}
	}
	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	return tokens
}

func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	var stack []rune

	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if isOperator(rune(token[0])) {
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[rune(token[0])] {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, rune(token[0]))
		} else if token == "(" {
			stack = append(stack, '(')
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				output = append(output, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errors.New("несоответствующая скобка")
			}
			stack = stack[:len(stack)-1]
		} else {
			return nil, fmt.Errorf("неизвестный токен: %s", token)
		}
	}

	for len(stack) > 0 {
		output = append(output, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if isNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else if isOperator(rune(token[0])) {
			if len(stack) < 2 {
				return 0, errors.New("недостаточно операндов")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch rune(token[0]) {
			case '+':
				stack = append(stack, a+b)
			case '-':
				stack = append(stack, a-b)
			case '*':
				stack = append(stack, a*b)
			case '/':
				if b == 0 {
					return 0, errors.New("деление на ноль")
				}
				stack = append(stack, a/b)
			default:
				return 0, fmt.Errorf("неизвестный оператор: %s", token)
			}
		} else {
			return 0, fmt.Errorf("неизвестный токен: %s", token)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("ошибка в вычислении выражения")
	}

	return stack[0], nil
}

func isNumber(token string) bool {
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return true
	}
	return false
}

func isOperator(char rune) bool {
	switch char {
	case '+', '-', '*', '/':
		return true
	default:
		return false
	}
}

func main() {
	result, err := Calc("3 + 2")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("Результат:", result)
}
