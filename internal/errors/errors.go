package errors

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

// BadRequestError represents a client-side input error.
type BadRequestError struct {
	Msg string
}

func (e *BadRequestError) Error() string {
	return e.Msg
}

func NewBadRequestError(msg string) *BadRequestError {
	return &BadRequestError{Msg: msg}
}

// ConvertError maps application errors to appropriate HTTP responses.
func ConvertError(w http.ResponseWriter, err error) {
	var status int

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
	case errors.As(err, new(*BadRequestError)):
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
	}

	http.Error(w, err.Error(), status)
}
