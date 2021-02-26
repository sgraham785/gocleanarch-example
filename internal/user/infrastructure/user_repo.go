package infrastructure

import (
	"github.com/sgraham785/gocleanarch-example/internal/user/entity"
)

//go:generate mockgen -destination=../mock/user_repo_mock.go -package=mock github.com/sgraham785/gocleanarch-example/internal/user/infrastructure Reader,Writer,UserRepo

// Reader interface
type Reader interface {
	Get(id entity.ID) (*entity.User, error)
	Search(query string) ([]*entity.User, error)
	List() ([]*entity.User, error)
}

// Writer user writer
type Writer interface {
	Create(e *entity.User) (entity.ID, error)
	Update(e *entity.User) error
	Delete(id entity.ID) error
}

// UserRepo interface
type UserRepo interface {
	Reader
	Writer
}
