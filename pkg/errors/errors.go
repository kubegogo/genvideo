package errors

import "fmt"

var (
	ErrNotFound       = fmt.Errorf("resource not found")
	ErrInvalidInput   = fmt.Errorf("invalid input")
	ErrUnauthorized   = fmt.Errorf("unauthorized")
	ErrInternalServer = fmt.Errorf("internal server error")
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
