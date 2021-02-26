package adapter

import (
	"fmt"
	"net/http"

	"github.com/rs/xid"
	bookEntity "github.com/sgraham785/gocleanarch-example/internal/book/entity"
	bookUseCase "github.com/sgraham785/gocleanarch-example/internal/book/usecase"
	"github.com/sgraham785/gocleanarch-example/internal/borrow/usecase"
	userEntity "github.com/sgraham785/gocleanarch-example/internal/user/entity"
	userUseCase "github.com/sgraham785/gocleanarch-example/internal/user/usecase"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/sgraham785/gocleanarch-example/internal/borrow/entity"
)

// BorrowBookHTTP handler
func BorrowBookHTTP(bookUseCase bookUseCase.BookUseCase, userUseCase userUseCase.UserUseCase, borrowUseCase usecase.BorrowUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error borrowing book"
		vars := mux.Vars(r)
		bID, err := xid.FromString(vars["book_id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		b, err := bookUseCase.GetBook(bID)
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
		uID, err := xid.FromString(vars["user_id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		u, err := userUseCase.GetUser(uID)
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
	})
}

// ReturnBookHTTP handler
func ReturnBookHTTP(bookUseCase bookUseCase.BookUseCase, borrowUseCase usecase.BorrowUseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error returning book"
		vars := mux.Vars(r)
		bID, err := xid.FromString(vars["book_id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		b, err := bookUseCase.GetBook(bID)
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
	})
}

// BorrowRouter make url handlers
func BorrowRouter(r *mux.Router, n negroni.Negroni, bookUseCase bookUseCase.BookUseCase, userUseCase userUseCase.UserUseCase, borrowUseCase usecase.BorrowUseCase) {
	r.Handle("/v1/loan/borrow/{book_id}/{user_id}", n.With(
		negroni.Wrap(BorrowBookHTTP(bookUseCase, userUseCase, borrowUseCase)),
	)).Methods("GET", "OPTIONS").Name("borrowBook")

	r.Handle("/v1/loan/return/{book_id}", n.With(
		negroni.Wrap(ReturnBookHTTP(bookUseCase, borrowUseCase)),
	)).Methods("GET", "OPTIONS").Name("returnBook")
}
