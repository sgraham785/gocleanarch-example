package usecase

import (
	"strings"
	"time"

	"github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/sgraham785/gocleanarch-example/internal/book/infrastructure"
	"github.com/sgraham785/gocleanarch-example/pkg/logger"
	"github.com/sgraham785/gocleanarch-example/pkg/server"
)

//go:generate mockgen -destination=../mock/book_usecase_mock.go -package=mock github.com/sgraham785/gocleanarch-example/internal/book/usecase BookUseCase

// BookUseCase is the interface that provides the methods.
type BookUseCase interface {
	GetBook(id string) (*entity.Book, error)
	SearchBooks(query string) ([]*entity.Book, error)
	ListBooks() ([]*entity.Book, error)
	CreateBook(title string, author string, pages int, quantity int) (entity.ID, error)
	UpdateBook(e *entity.Book) error
	DeleteBook(id string) error
}

type bookUseCase struct {
	repo infrastructure.BookRepo
	log  *logger.Logger
}

// New create new book usecase
func New(s *server.Server, r infrastructure.BookRepo) BookUseCase {
	return &bookUseCase{
		repo: r,
		log:  s.Log,
	}
}

// CreateBook create a book
func (u *bookUseCase) CreateBook(title string, author string, pages int, quantity int) (entity.ID, error) {
	b, err := entity.New(title, author, pages, quantity)
	if err != nil {
		return b.ID, err
	}
	return u.repo.Create(b)
}

// GetBook get a book
func (u *bookUseCase) GetBook(id string) (*entity.Book, error) {
	bID, err := entity.IDFromString(id)
	b, err := u.repo.Get(bID)
	if b == nil {
		return nil, entity.ErrBookNotFound
	}
	if err != nil {
		return nil, err
	}

	return b, nil
}

// SearchBooks search books
func (u *bookUseCase) SearchBooks(query string) ([]*entity.Book, error) {
	books, err := u.repo.Search(strings.ToLower(query))
	if err != nil {
		return nil, err
	}
	if len(books) == 0 {
		return nil, entity.ErrBookNotFound
	}
	return books, nil
}

// ListBooks list books
func (u *bookUseCase) ListBooks() ([]*entity.Book, error) {
	books, err := u.repo.List()
	if err != nil {
		return nil, err
	}
	if len(books) == 0 {
		return nil, entity.ErrBookNotFound
	}
	return books, nil
}

// DeleteBook Delete a book
func (u *bookUseCase) DeleteBook(id string) error {
	_, err := u.GetBook(id)
	if err != nil {
		return err
	}

	bID, _ := entity.IDFromString(id)
	return u.repo.Delete(bID)
}

// UpdateBook Update a book
func (u *bookUseCase) UpdateBook(e *entity.Book) error {
	err := e.Validate()
	if err != nil {
		return err
	}
	e.UpdatedAt = time.Now()
	return u.repo.Update(e)
}
