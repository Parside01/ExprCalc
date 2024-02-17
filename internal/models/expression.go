package models

import (
	"encoding/json"
	"errors"
	"strings"
	"unicode"

	"github.com/twharmon/gouid"
)

type Expression struct {
	Result            int    `json:"result"`
	Expression        string `json:"expression"`
	GUID              string `json:"guid"`
	ExecuteTime       int64  `json:"execute-time"`
	ExpectExucuteTime int64  `json:"expect-execute-time"`
	Err               error  `json:"err"`
	IsDone            bool   `json:"is-done"`
	WorkerID          string `json:"worker-id"`
}

func NewExpression(exp string) *Expression {
	guid := gouid.Bytes(16)
	return &Expression{Expression: exp, GUID: guid.String(), IsDone: false}
}

func (e *Expression) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Expression) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &e)
}

var (
	ErrUnequalNumStaples = errors.New("unequal number of opposite staples")
	ErrDevidedByZero     = errors.New("divided by zero")
	ErrInvalidChar       = errors.New("invalid char")
)

func (e *Expression) IsValidMathExpression() error {
	stack := []rune{}
	expression := e.Expression

	if strings.Contains(expression, "/0") || strings.Contains(expression, "/ 0") {
		return ErrDevidedByZero
	}

	for _, char := range expression {
		if unicode.IsDigit(char) || unicode.IsSpace(char) {
			continue
		} else if char == '(' {
			stack = append(stack, '(')
		} else if char == ')' {
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return ErrUnequalNumStaples
			}
			stack = stack[:len(stack)-1]
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			if len(stack) == 0 || stack[len(stack)-1] == '(' {
				return ErrUnequalNumStaples
			}
		} else {
			return ErrInvalidChar
		}
	}

	return nil
}
