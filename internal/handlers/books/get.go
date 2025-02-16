package books

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/mikhailsoldatkin/book_store/internal/models"
	"github.com/mikhailsoldatkin/book_store/internal/request"
)

// Get fetches a Book from database by ID.
func Get(db *gorm.DB, r *http.Request) (any, error) {
	bookID, err := request.ParseUintParam(r, "id")
	if err != nil {
		return nil, err
	}

	var book models.Book
	if err = db.First(&book, bookID).Error; err != nil {
		return nil, err
	}

	return book, nil
}
