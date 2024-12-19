package calculation

import "errors"

var (
	ErrInvalidExpression   = errors.New("invalid expression")
	ErrDivisionByZero      = errors.New("division by zero")
	ErrUnknowOperator      = errors.New("unknow operator")
	ErrInvalidCharacter    = errors.New("invalid character")
	ErrInvalidToken        = errors.New("invalid token")
	ErrMismatchParentheses = errors.New("mismatched parentheses")
)
