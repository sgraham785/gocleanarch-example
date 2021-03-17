package usecase

import (
	bookEntity "github.com/sgraham785/gocleanarch-example/internal/book/entity"
	bookUseCase "github.com/sgraham785/gocleanarch-example/internal/book/usecase"
	"github.com/sgraham785/gocleanarch-example/internal/borrow/entity"
	userEntity "github.com/sgraham785/gocleanarch-example/internal/user/entity"
	userUseCase "github.com/sgraham785/gocleanarch-example/internal/user/usecase"
	"github.com/sgraham785/gocleanarch-example/pkg/logger"
	"github.com/sgraham785/gocleanarch-example/pkg/server"
)

//go:generate mockgen -destination=../mock/borrow_usecase_mock.go -package=mock github.com/sgraham785/gocleanarch-example/internal/borrow/usecase BorrowUseCase

// BorrowUseCase is the interface that provides the methods.
type BorrowUseCase interface {
	Borrow(u *userEntity.User, b *bookEntity.Book) error
	Return(b *bookEntity.Book) error
}

type borrowUseCase struct {
	userUseCase userUseCase.UserUseCase
	bookUseCase bookUseCase.BookUseCase
	log         *logger.Logger
}

// New create new borrow use case
func New(s *server.Server, u userUseCase.UserUseCase, b bookUseCase.BookUseCase) BorrowUseCase {
	return &borrowUseCase{
		userUseCase: u,
		bookUseCase: b,
		log:         s.Log,
	}
}

// Borrow borrow a book to an user
func (s *borrowUseCase) Borrow(u *userEntity.User, b *bookEntity.Book) error {
	u, err := s.userUseCase.GetUser(u.ID.String())
	if err != nil {
		return err
	}
	b, err = s.bookUseCase.GetBook(b.ID.String())
	if err != nil {
		return err
	}
	if b.Quantity <= 0 {
		return entity.ErrNotEnoughBooks
	}

	err = u.AddBook(b.ID)
	if err != nil {
		return err
	}
	err = s.userUseCase.UpdateUser(u)
	if err != nil {
		return err
	}
	b.Quantity--
	err = s.bookUseCase.UpdateBook(b)
	if err != nil {
		return err
	}
	return nil
}

//Return return a book
func (s *borrowUseCase) Return(b *bookEntity.Book) error {
	b, err := s.bookUseCase.GetBook(b.ID.String())
	if err != nil {
		return err
	}

	all, err := s.userUseCase.ListUsers()
	if err != nil {
		return err
	}
	borrowed := false
	var borrowedBy userEntity.ID
	for _, u := range all {
		_, err := u.GetBook(b.ID)
		if err != nil {
			continue
		}
		borrowed = true
		borrowedBy = u.ID
		break
	}
	if !borrowed {
		return entity.ErrBookNotBorrowed
	}
	u, err := s.userUseCase.GetUser(borrowedBy.String())
	if err != nil {
		return err
	}
	err = u.RemoveBook(b.ID)
	if err != nil {
		return err
	}
	err = s.userUseCase.UpdateUser(u)
	if err != nil {
		return err
	}
	b.Quantity++
	err = s.bookUseCase.UpdateBook(b)
	if err != nil {
		return err
	}

	return nil
}
