package common

import (
	"net/http"
)

type Error struct {
	StatusCode int
	Message    []string
	Err        error
}

// NewError /**
func NewError(c int, e ...error) Error {
	if len(e) > 0 {
		return Error{StatusCode: c, Err: e[0]}
	}

	return Error{StatusCode: c}
}

// PanicInternalServerError /**
func PanicInternalServerError(err error) {
	panic(NewError(http.StatusInternalServerError, err))
}

// PanicBadRequest /**
func PanicBadRequest(err error) {
	panic(NewError(http.StatusBadRequest, err))
}

// PanicTooManyRequests /**
func PanicTooManyRequests() {
	panic(NewError(http.StatusTooManyRequests))
}

// PanicNotFound /**
func PanicNotFound() {
	panic(NewError(http.StatusNotFound))
}

// PanicUnauthorized /**
func PanicUnauthorized() {
	panic(NewError(http.StatusUnauthorized))
}
