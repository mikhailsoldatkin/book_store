package request

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mikhailsoldatkin/book_store/internal/errors"
)

// ParseUintParam ...
func ParseUintParam(r *http.Request, paramName string) (uint, error) {
	params := mux.Vars(r)

	param, err := strconv.Atoi(params[paramName])
	if err != nil {
		return 0, errors.NewBadRequestError("can not parse param " + paramName)
	}

	return uint(param), nil
}
