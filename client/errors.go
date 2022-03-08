package client

import "fmt"

type RequestError struct {
	Message    string
	StatusCode int
}

type NotFoundError struct {
	Message    string
	StatusCode int
}

func NewRequestError(StatusCode int, message string) RequestError {
	return RequestError{
		Message:    message,
		StatusCode: StatusCode,
	}
}

func (err RequestError) Error() string {
	return fmt.Sprintf("(%d) %s", err.StatusCode, err.Message)
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{
		Message:    message,
		StatusCode: 404,
	}
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("NotFoundError (%d) %s", err.StatusCode, err.Message)
}
