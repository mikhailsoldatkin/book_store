package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	"github.com/mikhailsoldatkin/book_store/internal/errors"
)

// HandlerFunc is a function type that processes requests with a database connection.
type HandlerFunc func(db *gorm.DB, r *http.Request) (any, error)

// WrapHandler converts a HandlerFunc into a http.HandlerFunc.
func WrapHandler(db *gorm.DB, handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := handler(db, r)
		if err != nil {
			errors.ConvertError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
