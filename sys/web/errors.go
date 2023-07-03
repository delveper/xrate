package web

import (
	"errors"
	"net/http"
)

var (
	ErrClientError  = errors.New("client error")
	ErrServerError  = errors.New("server error")
	ErrUnknownError = errors.New("unknown error")
)

// RequestError type is an error that is used to indicate that a request failed.
type RequestError struct {
	Value      error
	StatusCode int
}

// NewRequestError creates a new request error.
func NewRequestError(err error, code int) error {
	return &RequestError{Value: err, StatusCode: code}
}

func (e *RequestError) Error() string {
	return e.Value.Error()
}

// ErrorResponse type is an error that is used to indicate that a response failed.
type ErrorResponse struct {
	Error   string         `json:"error"`
	Details map[string]any `json:"details,omitempty"`
}

// ShutdownError type is an error that is used to indicate that the web application is shutting down.
type ShutdownError struct {
	Message string
}

// NewShutdownError creates a new shutdown error.
func NewShutdownError(msg string) error {
	return &ShutdownError{Message: msg}
}

func (e *ShutdownError) Error() string {
	return e.Message
}

// IsError checks if the error is a known type returning its value if true.
func IsError[E error](err error) (errVal E, ok bool) {
	return errVal, errors.As(err, &errVal)
}

// ErrFromStatusCode categorizes HTTP response status codes and returns corresponding errors.
func ErrFromStatusCode(code int) error {
	switch code {
	case http.StatusOK, http.StatusAccepted:
		return nil

	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden:
		return NewRequestError(ErrClientError, code)

	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return NewRequestError(ErrServerError, code)

	default:
		return NewRequestError(ErrUnknownError, code)
	}
}
