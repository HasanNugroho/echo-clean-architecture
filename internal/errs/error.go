// internal/errs/error.go
package errs

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	Code    int
	Message string
	Err     error
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *CustomError) Unwrap() error {
	return e.Err
}

func (e *CustomError) StatusCode() int {
	return e.Code
}

func (e *CustomError) MessageText() string {
	return e.Message
}

func BadRequest(msg string, err error) *CustomError {
	return &CustomError{Code: http.StatusBadRequest, Message: msg, Err: err}
}

func NotFound(msg string, err error) *CustomError {
	return &CustomError{Code: http.StatusNotFound, Message: msg, Err: err}
}

func Internal(msg string, err error) *CustomError {
	return &CustomError{Code: http.StatusInternalServerError, Message: msg, Err: err}
}
