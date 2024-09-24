package http

import (
	"fmt"
	"net/http"
)

func badRequestError(format string, v ...any) error {
	if len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	return statusError{http.StatusBadRequest, format}
}

func unauthorizedError(format string, v ...any) error {
	if len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	return statusError{http.StatusUnauthorized, format}
}

func requestTimeoutError(format string, v ...any) error {
	if len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	return statusError{http.StatusRequestTimeout, format}
}

func requestEntityTooLargeError(format string, v ...any) error {
	if len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	return statusError{http.StatusRequestEntityTooLarge, format}
}

func unsupportedMediaTypeError(format string, v ...any) error {
	if len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	return statusError{http.StatusUnsupportedMediaType, format}
}

// statusError is an error used to respond to a request with an HTTP status.
type statusError struct {
	code int
	text string
}

func (e statusError) Error() string { return http.StatusText(e.code) + ": " + e.text }

func Error(w http.ResponseWriter, r *http.Request, err error) {
	if se, ok := err.(statusError); ok {
		http.Error(w, se.text, se.code)
		return
	}
	// do some stuff with domain specific errors
	http.Error(w, http.StatusText(500), http.StatusInternalServerError)
}
