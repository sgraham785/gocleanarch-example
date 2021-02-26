package adapter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	bookEntity "github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/sgraham785/gocleanarch-example/internal/borrow/adapter"
	userEntity "github.com/sgraham785/gocleanarch-example/internal/user/entity"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	bookMock "github.com/sgraham785/gocleanarch-example/internal/book/mock"
	borrowMock "github.com/sgraham785/gocleanarch-example/internal/borrow/mock"
	userMock "github.com/sgraham785/gocleanarch-example/internal/user/mock"
	"github.com/stretchr/testify/assert"
)

func Test_BorrowBookHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	userMock := userMock.NewMockUserUseCase(controller)
	bookMock := bookMock.NewMockBookUseCase(controller)
	borrowMock := borrowMock.NewMockBorrowUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	adapter.BorrowRouter(r, *n, bookMock, userMock, borrowMock)
	path, err := r.GetRoute("borrowBook").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/loan/borrow/{book_id}/{user_id}", path)
	handler := adapter.BorrowBookHTTP(bookMock, userMock, borrowMock)
	r.Handle("/v1/loan/borrow/{book_id}/{user_id}", handler)
	t.Run("book not found", func(t *testing.T) {
		bID := bookEntity.NewID()
		uID := userEntity.NewID()
		bookMock.EXPECT().GetBook(bID).Return(nil, bookEntity.ErrBookNotFound)
		ts := httptest.NewServer(r)
		defer ts.Close()
		res, err := http.Get(fmt.Sprintf("%s/v1/loan/borrow/%s/%s", ts.URL, bID.String(), uID.String()))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
	t.Run("user not found", func(t *testing.T) {
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		uID := userEntity.NewID()
		bookMock.EXPECT().GetBook(b.ID).Return(b, nil)
		userMock.EXPECT().GetUser(uID).Return(nil, userEntity.ErrUserNotFound)
		ts := httptest.NewServer(r)
		defer ts.Close()
		res, err := http.Get(fmt.Sprintf("%s/v1/loan/borrow/%s/%s", ts.URL, b.ID.String(), uID.String()))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
	t.Run("success", func(t *testing.T) {
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		u := &userEntity.User{
			ID: userEntity.NewID(),
		}
		bookMock.EXPECT().GetBook(b.ID).Return(b, nil)
		userMock.EXPECT().GetUser(u.ID).Return(u, nil)
		borrowMock.EXPECT().Borrow(u, b).Return(nil)
		ts := httptest.NewServer(r)
		defer ts.Close()
		res, err := http.Get(fmt.Sprintf("%s/v1/loan/borrow/%s/%s", ts.URL, b.ID.String(), u.ID.String()))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}

func Test_ReturnBookHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	userMock := userMock.NewMockUserUseCase(controller)
	bookMock := bookMock.NewMockBookUseCase(controller)
	borrowMock := borrowMock.NewMockBorrowUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	adapter.BorrowRouter(r, *n, bookMock, userMock, borrowMock)
	path, err := r.GetRoute("returnBook").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/loan/return/{book_id}", path)
	handler := adapter.ReturnBookHTTP(bookMock, borrowMock)
	r.Handle("/v1/loan/return/{book_id}", handler)
	t.Run("book not found", func(t *testing.T) {
		bID := bookEntity.NewID()
		bookMock.EXPECT().GetBook(bID).Return(nil, bookEntity.ErrBookNotFound)
		ts := httptest.NewServer(r)
		defer ts.Close()
		res, err := http.Get(fmt.Sprintf("%s/v1/loan/return/%s", ts.URL, bID.String()))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
	t.Run("success", func(t *testing.T) {
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		bookMock.EXPECT().GetBook(b.ID).Return(b, nil)
		borrowMock.EXPECT().Return(b).Return(nil)
		ts := httptest.NewServer(r)
		defer ts.Close()
		res, err := http.Get(fmt.Sprintf("%s/v1/loan/return/%s", ts.URL, b.ID.String()))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}
