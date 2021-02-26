package infrastructure

import (
	"github.com/sgraham785/gocleanarch-example/internal/book/entity"
)

//go:generate mockgen -destination=../mock/book_repo_mock.go -package=mock github.com/sgraham785/gocleanarch-example/internal/book/infrastructure Reader,Writer,BookRepo

// Reader interface
type Reader interface {
	Get(id entity.ID) (*entity.Book, error)
	Search(query string) ([]*entity.Book, error)
	List() ([]*entity.Book, error)
}

// Writer interface
type Writer interface {
	Create(e *entity.Book) (entity.ID, error)
	Update(e *entity.Book) error
	Delete(id entity.ID) error
}

// BookRepo interface
type BookRepo interface {
	Reader
	Writer
}
