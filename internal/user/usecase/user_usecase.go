package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/sgraham785/gocleanarch-example/internal/user/entity"
	"github.com/sgraham785/gocleanarch-example/internal/user/infrastructure"
)

//go:generate mockgen -destination=../mock/user_usecase_mock.go -package=mock github.com/sgraham785/gocleanarch-example/internal/user/usecase UserUseCase

// UserUseCase is the interface that provides the methods.
type UserUseCase interface {
	GetUser(id entity.ID) (*entity.User, error)
	SearchUsers(query string) ([]*entity.User, error)
	ListUsers() ([]*entity.User, error)
	CreateUser(email, password, firstName, lastName string) (entity.ID, error)
	UpdateUser(e *entity.User) error
	DeleteUser(id entity.ID) error
}

type userUseCase struct {
	repo infrastructure.UserRepo
}

// New create new user use case
func New(r infrastructure.UserRepo) UserUseCase {
	return &userUseCase{
		repo: r,
	}
}

// CreateUser create an user
func (s *userUseCase) CreateUser(email, password, firstName, lastName string) (entity.ID, error) {
	fmt.Println("CreateUser")
	e, err := entity.New(email, password, firstName, lastName)
	if err != nil {
		return e.ID, err
	}
	return s.repo.Create(e)
}

// GetUser gets an user
func (s *userUseCase) GetUser(id entity.ID) (*entity.User, error) {
	return s.repo.Get(id)
}

// SearchUsers searches users
func (s *userUseCase) SearchUsers(query string) ([]*entity.User, error) {
	return s.repo.Search(strings.ToLower(query))
}

// ListUsers lists users
func (s *userUseCase) ListUsers() ([]*entity.User, error) {
	return s.repo.List()
}

// DeleteUser deletes an user
func (s *userUseCase) DeleteUser(id entity.ID) error {
	u, err := s.GetUser(id)
	if u == nil {
		return entity.ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if len(u.Books) > 0 {
		return entity.ErrUserCannotBeDeleted
	}
	return s.repo.Delete(id)
}

// UpdateUser updates an user
func (s *userUseCase) UpdateUser(e *entity.User) error {
	err := e.Validate()
	if err != nil {
		return entity.ErrInvalidUserEntity
	}
	e.UpdatedAt = time.Now()
	return s.repo.Update(e)
}
