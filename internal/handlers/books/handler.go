package books

import (
	"gorm.io/gorm"
)

// Handler represents a handler for HTTP requests.
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new instance of Handler.
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}
