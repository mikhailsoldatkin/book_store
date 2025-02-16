package handlers

import (
	"github.com/mikhailsoldatkin/book_store/internal/handlers/books"

	"net/http"
)

// Route describes route.
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler HandlerFunc
}

var Routes = []Route{
	{
		Name:    "Get Book",
		Method:  http.MethodGet,
		Pattern: "/books/{id:[0-9]+}",
		Handler: books.Get,
	},
}
