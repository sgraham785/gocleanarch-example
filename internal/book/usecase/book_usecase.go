package usecase

import (
	"strings"
	"time"

	"github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/sgraham785/gocleanarch-example/internal/book/infrastructure"
)

//go:generate mockgen -destination=../mock/book_usecase_mock.go -package=mock github.com/sgraham785/gocleanarch-example/internal/book/usecase BookUseCase

// BookUseCase is the interface that provides the methods.
type BookUseCase interface {
	GetBook(id entity.ID) (*entity.Book, error)
	SearchBooks(query string) ([]*entity.Book, error)
	ListBooks() ([]*entity.Book, error)
	CreateBook(title string, author string, pages int, quantity int) (entity.ID, error)
	UpdateBook(e *entity.Book) error
	DeleteBook(id entity.ID) error
}

type bookUseCase struct {
	repo infrastructure.BookRepo
}

// New create new book usecase
func New(r infrastructure.BookRepo) BookUseCase {
	return &bookUseCase{
		repo: r,
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
func (u *bookUseCase) GetBook(id entity.ID) (*entity.Book, error) {
	b, err := u.repo.Get(id)
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
func (u *bookUseCase) DeleteBook(id entity.ID) error {
	_, err := u.GetBook(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(id)
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
