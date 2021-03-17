package adapter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/sgraham785/gocleanarch-example/internal/book/usecase"
	"github.com/sgraham785/gocleanarch-example/pkg/server"
)

// BookHTTP JSON data
type BookHTTP struct {
	ID       entity.ID `json:"id"`
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	Pages    int       `json:"pages"`
	Quantity int       `json:"quantity"`
}

// ListBooksHTTP handler
func ListBooksHTTP(u usecase.BookUseCase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading books"
		var data []*entity.Book
		var err error
		title := r.URL.Query().Get("title")
		switch {
		case title == "":
			data, err = u.ListBooks()
		default:
			data, err = u.SearchBooks(title)
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrBookNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		var toJ []*BookHTTP
		for _, d := range data {
			toJ = append(toJ, &BookHTTP{
				ID:       d.ID,
				Title:    d.Title,
				Author:   d.Author,
				Pages:    d.Pages,
				Quantity: d.Quantity,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

// CreateBookHTTP handler
func CreateBookHTTP(u usecase.BookUseCase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding book"
		var input struct {
			Title    string `json:"title"`
			Author   string `json:"author"`
			Pages    int    `json:"pages"`
			Quantity int    `json:"quantity"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		id, err := u.CreateBook(input.Title, input.Author, input.Pages, input.Quantity)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		toJ := &BookHTTP{
			ID:       id,
			Title:    input.Title,
			Author:   input.Author,
			Pages:    input.Pages,
			Quantity: input.Quantity,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

// GetBookHTTP handler
func GetBookHTTP(u usecase.BookUseCase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading book"

		if bookID := chi.URLParam(r, "bookID"); bookID != "" {
			data, err := u.GetBook(bookID)
			if err != nil && err != entity.ErrBookNotFound {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorMessage))
				return
			}
			if data == nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errorMessage))
				return
			}
			toJ := &BookHTTP{
				ID:       data.ID,
				Title:    data.Title,
				Author:   data.Author,
				Pages:    data.Pages,
				Quantity: data.Quantity,
			}
			if err := json.NewEncoder(w).Encode(toJ); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorMessage))
			}

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

	})
}

// DeleteBookHTTP handler
func DeleteBookHTTP(u usecase.BookUseCase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing bookmark"
		if bookID := chi.URLParam(r, "bookID"); bookID != "" {
			err := u.DeleteBook(bookID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorMessage))
				return
			}
		}
	})
}

// HTTPRoutes defines http routes for books
func HTTPRoutes(s *server.Server, u usecase.BookUseCase) {
	// RESTy routes for "books" resource
	s.Router.Chi.Route("/book", func(r chi.Router) {
		// r.With(paginate).Get("/", ListBooksHTTP(u))
		r.Get("/", ListBooksHTTP(u))
		r.Post("/", CreateBookHTTP(u))     // POST /book
		r.Get("/search", ListBooksHTTP(u)) // GET /book/search?title=something

		r.Route("/{bookID}", func(r chi.Router) {
			r.Get("/", GetBookHTTP(u)) // GET /book/123
			// r.Put("/", UpdateArticle)    // PUT /book/123
			r.Delete("/", DeleteBookHTTP(u)) // DELETE /book/123
		})

		// GET /book/whats-up
		// r.With(BookCtx).Get("/{bookTitle:[a-z-]+}", GetBookHTTP(u))
	})
}
