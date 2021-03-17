package adapter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	bookEntity "github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/sgraham785/gocleanarch-example/internal/borrow/adapter"
	userEntity "github.com/sgraham785/gocleanarch-example/internal/user/entity"
	"github.com/sgraham785/gocleanarch-example/pkg/router"
	"github.com/sgraham785/gocleanarch-example/pkg/server"

	"github.com/golang/mock/gomock"
	bookMock "github.com/sgraham785/gocleanarch-example/internal/book/mock"
	borrowMock "github.com/sgraham785/gocleanarch-example/internal/borrow/mock"
	userMock "github.com/sgraham785/gocleanarch-example/internal/user/mock"
	"github.com/stretchr/testify/assert"
)

func TestBorrowBookHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	userMock := userMock.NewMockUserUseCase(controller)
	bookMock := bookMock.NewMockBookUseCase(controller)
	borrowMock := borrowMock.NewMockBorrowUseCase(controller)
	r := router.NewChiRouter()
	s := &server.Server{
		Router: r,
	}
	adapter.HTTPRoutes(s, bookMock, userMock, borrowMock)
	handler := adapter.BorrowBookHTTP(bookMock, userMock, borrowMock)
	r.Chi.Handle("/borrow/{bookID}/{userID}", handler)

	t.Run("book not found", func(t *testing.T) {
		bID := bookEntity.NewID()
		uID := userEntity.NewID()
		bookMock.EXPECT().GetBook(bID.String()).Return(nil, bookEntity.ErrBookNotFound)
		ts := httptest.NewServer(r.Chi)
		defer ts.Close()
		res, err := http.Get(fmt.Sprintf("%s/borrow/%s/%s", ts.URL, bID.String(), uID.String()))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("user not found", func(t *testing.T) {
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		uID := userEntity.NewID()
		bookMock.EXPECT().GetBook(b.ID.String()).Return(b, nil)
		userMock.EXPECT().GetUser(uID.String()).Return(nil, userEntity.ErrUserNotFound)
		ts := httptest.NewServer(r.Chi)
		defer ts.Close()
		res, err := http.Get(fmt.Sprintf("%s/borrow/%s/%s", ts.URL, b.ID.String(), uID.String()))
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
		bookMock.EXPECT().GetBook(b.ID.String()).Return(b, nil)
		userMock.EXPECT().GetUser(u.ID.String()).Return(u, nil)
		borrowMock.EXPECT().Borrow(u, b).Return(nil)
		ts := httptest.NewServer(r.Chi)
		defer ts.Close()
		res, err := http.Get(fmt.Sprintf("%s/borrow/%s/%s", ts.URL, b.ID.String(), u.ID.String()))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}

func TestReturnBookHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	userMock := userMock.NewMockUserUseCase(controller)
	bookMock := bookMock.NewMockBookUseCase(controller)
	borrowMock := borrowMock.NewMockBorrowUseCase(controller)
	r := router.NewChiRouter()
	s := &server.Server{
		Router: r,
	}
	adapter.HTTPRoutes(s, bookMock, userMock, borrowMock)
	h := adapter.ReturnBookHTTP(bookMock, borrowMock)
	r.Chi.Handle("/return/{bookID}", h)

	t.Run("book not found", func(t *testing.T) {
		bID := bookEntity.NewID()
		bookMock.EXPECT().GetBook(bID.String()).Return(nil, bookEntity.ErrBookNotFound)
		ts := httptest.NewServer(r.Chi)
		defer ts.Close()
		res, err := http.Get(fmt.Sprintf("%s/return/%s", ts.URL, bID.String()))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		bookMock.EXPECT().GetBook(b.ID.String()).Return(b, nil)
		borrowMock.EXPECT().Return(b).Return(nil)
		ts := httptest.NewServer(r.Chi)
		defer ts.Close()
		res, err := http.Get(fmt.Sprintf("%s/return/%s", ts.URL, b.ID.String()))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}
