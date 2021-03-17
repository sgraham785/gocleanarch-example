package adapter

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	bookEntity "github.com/sgraham785/gocleanarch-example/internal/book/entity"
	bookUseCase "github.com/sgraham785/gocleanarch-example/internal/book/usecase"
	"github.com/sgraham785/gocleanarch-example/internal/borrow/usecase"
	userEntity "github.com/sgraham785/gocleanarch-example/internal/user/entity"
	userUseCase "github.com/sgraham785/gocleanarch-example/internal/user/usecase"
	"github.com/sgraham785/gocleanarch-example/pkg/server"

	"github.com/sgraham785/gocleanarch-example/internal/borrow/entity"
)

// BorrowBookHTTP handler
func BorrowBookHTTP(bookUseCase bookUseCase.BookUseCase, userUseCase userUseCase.UserUseCase, borrowUseCase usecase.BorrowUseCase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error borrowing book"
		if bookID := chi.URLParam(r, "bookID"); bookID != "" {
			bID, err := bookEntity.IDFromString(bookID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorMessage))
				return
			}
			b, err := bookUseCase.GetBook(bID.String())
			if err != nil && err != bookEntity.ErrBookNotFound {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorMessage))
				return
			}
			if b == nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errorMessage))
				return
			}
			if userID := chi.URLParam(r, "userID"); userID != "" {
				uID, err := userEntity.IDFromString(userID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(errorMessage))
					return
				}
				u, err := userUseCase.GetUser(uID.String())
				if err != nil && err != userEntity.ErrUserNotFound {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(errorMessage))
					return
				}
				if u == nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(errorMessage))
					return
				}
				err = borrowUseCase.Borrow(u, b)
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(errorMessage))
					return
				}
				w.WriteHeader(http.StatusCreated)
			}
		}
	})
}

// ReturnBookHTTP handler
func ReturnBookHTTP(bookUseCase bookUseCase.BookUseCase, borrowUseCase usecase.BorrowUseCase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error returning book"
		if bookID := chi.URLParam(r, "bookID"); bookID != "" {
			bID, err := bookEntity.IDFromString(bookID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorMessage))
				return
			}
			b, err := bookUseCase.GetBook(bID.String())
			if err != nil && err != bookEntity.ErrBookNotFound {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorMessage))
				return
			}
			if b == nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errorMessage))
				return
			}
			err = borrowUseCase.Return(b)
			if err != nil && err != entity.ErrBookNotBorrowed {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errorMessage))
				return
			}
			w.WriteHeader(http.StatusCreated)
		}
	})
}

// HTTPRoutes make url handlers
func HTTPRoutes(s *server.Server, bookUseCase bookUseCase.BookUseCase, userUseCase userUseCase.UserUseCase, borrowUseCase usecase.BorrowUseCase) {
	// RESTy routes for "books" resource
	s.Router.Chi.Route("/borrow", func(r chi.Router) {
		r.Post("/{bookID}/{userID}", BorrowBookHTTP(bookUseCase, userUseCase, borrowUseCase))
		r.Post("/return/{bookID}", ReturnBookHTTP(bookUseCase, borrowUseCase))
	})
}
