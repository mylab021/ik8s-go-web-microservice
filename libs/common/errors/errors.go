package errors

import "fmt"

type ErrorCode int

const (
	ErrCodeNotFound ErrorCode = iota + 1000
	ErrCodeInvalidInput
	// ...
)

type AppError struct {
	Code    ErrorCode
	Message string
	Op      string
	Err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}
