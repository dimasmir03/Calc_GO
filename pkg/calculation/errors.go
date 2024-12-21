package calculation

import "errors"

var (
	ErrInvalidExpression   = errors.New("invalid expression")
	ErrDivisionByZero      = errors.New("division by zero")
	ErrInvalidCharacter    = errors.New("invalid character")
	ErrMismatchParentheses = errors.New("mismatched parentheses")
)
