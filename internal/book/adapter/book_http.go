package adapter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/sgraham785/gocleanarch-example/internal/book/usecase"
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
func ListBooksHTTP(u usecase.BookUseCase) http.Handler {
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
func CreateBookHTTP(u usecase.BookUseCase) http.Handler {
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
			log.Println(err.Error())
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
func GetBookHTTP(u usecase.BookUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading book"
		vars := mux.Vars(r)
		id, err := xid.FromString(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		data, err := u.GetBook(id)
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
	})
}

// DeleteBookHTTP handler
func DeleteBookHTTP(u usecase.BookUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing bookmark"
		vars := mux.Vars(r)
		id, err := xid.FromString(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		err = u.DeleteBook(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

// BookRouter defined http routes for handlers
func BookRouter(r *mux.Router, n negroni.Negroni, u usecase.BookUseCase) {
	r.Handle("/v1/book", n.With(
		negroni.Wrap(ListBooksHTTP(u)),
	)).Methods("GET", "OPTIONS").Name("listBooks")

	r.Handle("/v1/book", n.With(
		negroni.Wrap(CreateBookHTTP(u)),
	)).Methods("POST", "OPTIONS").Name("createBook")

	r.Handle("/v1/book/{id}", n.With(
		negroni.Wrap(GetBookHTTP(u)),
	)).Methods("GET", "OPTIONS").Name("getBook")

	r.Handle("/v1/book/{id}", n.With(
		negroni.Wrap(DeleteBookHTTP(u)),
	)).Methods("DELETE", "OPTIONS").Name("deleteBook")
}
