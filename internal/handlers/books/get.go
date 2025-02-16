package books

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mikhailsoldatkin/book_store/internal/models"
)

// Get обрабатывает запрос на получение книги по ID.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var book models.Book

	vars := mux.Vars(r)
	bookIDStr := vars["id"]

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid ID: %s", bookIDStr), http.StatusBadRequest)
		return
	}

	if err := h.db.First(&book, bookID).Error; err != nil {
		http.Error(w, fmt.Sprintf("failed to get book: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode book: %s", err.Error()), http.StatusInternalServerError)
	}
}
