package books

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mikhailsoldatkin/book_store/internal/models"
)

// List handles the HTTP request to fetch a list of books from the database.
func (h *Handler) List(w http.ResponseWriter, _ *http.Request) {
	var books []*models.Book

	if err := h.db.
		Find(&books).Error; err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch books:%s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch books:%s", err.Error()), http.StatusInternalServerError)
	}
}
