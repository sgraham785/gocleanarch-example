package usecase_test

import (
	"testing"

	bookEntity "github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/sgraham785/gocleanarch-example/internal/borrow/entity"
	"github.com/sgraham785/gocleanarch-example/internal/borrow/usecase"
	userEntity "github.com/sgraham785/gocleanarch-example/internal/user/entity"

	"github.com/golang/mock/gomock"
	bookMock "github.com/sgraham785/gocleanarch-example/internal/book/mock"
	userMock "github.com/sgraham785/gocleanarch-example/internal/user/mock"
	"github.com/stretchr/testify/assert"
)

func Test_borrowUseCase_Borrow(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	uMock := userMock.NewMockUserUseCase(controller)
	bMock := bookMock.NewMockBookUseCase(controller)
	uc := usecase.New(uMock, bMock)
	t.Run("user not found", func(t *testing.T) {
		u := &userEntity.User{
			ID: userEntity.NewID(),
		}
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		uMock.EXPECT().GetUser(u.ID).Return(nil, userEntity.ErrUserNotFound)
		err := uc.Borrow(u, b)
		assert.Equal(t, userEntity.ErrUserNotFound, err)
	})
	t.Run("book not found", func(t *testing.T) {
		u := &userEntity.User{
			ID: userEntity.NewID(),
		}
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		uMock.EXPECT().GetUser(u.ID).Return(u, nil)
		bMock.EXPECT().GetBook(b.ID).Return(nil, bookEntity.ErrBookNotFound)
		err := uc.Borrow(u, b)
		assert.Equal(t, bookEntity.ErrBookNotFound, err)
	})
	t.Run("not enough books to borrow", func(t *testing.T) {
		u := &userEntity.User{
			ID: userEntity.NewID(),
		}
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		b.Quantity = 0
		uMock.EXPECT().GetUser(u.ID).Return(u, nil)
		bMock.EXPECT().GetBook(b.ID).Return(b, nil)
		err := uc.Borrow(u, b)
		assert.Equal(t, entity.ErrNotEnoughBooks, err)
	})
	t.Run("book already borrowed", func(t *testing.T) {
		u := &userEntity.User{
			ID: userEntity.NewID(),
		}
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		u.AddBook(b.ID)
		b.Quantity = 1
		uMock.EXPECT().GetUser(u.ID).Return(u, nil)
		bMock.EXPECT().GetBook(b.ID).Return(b, nil)
		err := uc.Borrow(u, b)
		assert.Equal(t, entity.ErrBookAlreadyBorrowed, err)
	})
	t.Run("sucess", func(t *testing.T) {
		u := &userEntity.User{
			ID: userEntity.NewID(),
		}
		b := &bookEntity.Book{
			ID:       bookEntity.NewID(),
			Quantity: 10,
		}
		uMock.EXPECT().GetUser(u.ID).Return(u, nil)
		bMock.EXPECT().GetBook(b.ID).Return(b, nil)
		uMock.EXPECT().UpdateUser(u).Return(nil)
		bMock.EXPECT().UpdateBook(b).Return(nil)
		err := uc.Borrow(u, b)
		assert.Nil(t, err)
	})
}

func Test_borrowUseCase_Return(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	uMock := userMock.NewMockUserUseCase(controller)
	bMock := bookMock.NewMockBookUseCase(controller)
	uc := usecase.New(uMock, bMock)
	t.Run("book not found", func(t *testing.T) {
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		bMock.EXPECT().GetBook(b.ID).Return(nil, bookEntity.ErrBookNotFound)
		err := uc.Return(b)
		assert.Equal(t, bookEntity.ErrBookNotFound, err)
	})
	t.Run("book not borrowed", func(t *testing.T) {
		u := &userEntity.User{
			ID: userEntity.NewID(),
		}
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		bMock.EXPECT().GetBook(b.ID).Return(b, nil)
		uMock.EXPECT().ListUsers().Return([]*userEntity.User{u}, nil)
		err := uc.Return(b)
		assert.Equal(t, entity.ErrBookNotBorrowed, err)
	})
	t.Run("success", func(t *testing.T) {
		u := &userEntity.User{
			ID: userEntity.NewID(),
		}
		b := &bookEntity.Book{
			ID: bookEntity.NewID(),
		}
		u.AddBook(b.ID)
		bMock.EXPECT().GetBook(b.ID).Return(b, nil)
		uMock.EXPECT().GetUser(u.ID).Return(u, nil)
		uMock.EXPECT().ListUsers().Return([]*userEntity.User{u}, nil)
		uMock.EXPECT().UpdateUser(u).Return(nil)
		bMock.EXPECT().UpdateBook(b).Return(nil)
		err := uc.Return(b)
		assert.Nil(t, err)
	})
}
