package usecase

import (
	"strings"
	"time"

	"github.com/sgraham785/gocleanarch-example/internal/user/entity"
	"github.com/sgraham785/gocleanarch-example/internal/user/infrastructure"
	"github.com/sgraham785/gocleanarch-example/pkg/logger"
	"github.com/sgraham785/gocleanarch-example/pkg/server"
)

//go:generate mockgen -destination=../mock/user_usecase_mock.go -package=mock github.com/sgraham785/gocleanarch-example/internal/user/usecase UserUseCase

// UserUseCase is the interface that provides the methods.
type UserUseCase interface {
	GetUser(id string) (*entity.User, error)
	SearchUsers(query string) ([]*entity.User, error)
	ListUsers() ([]*entity.User, error)
	CreateUser(email, password, firstName, lastName string) (entity.ID, error)
	UpdateUser(e *entity.User) error
	DeleteUser(id string) error
}

type userUseCase struct {
	repo infrastructure.UserRepo
	log  *logger.Logger
}

// New create new user use case
func New(s *server.Server, r infrastructure.UserRepo) UserUseCase {
	return &userUseCase{
		repo: r,
		log:  s.Log,
	}
}

// CreateUser create an user
func (s *userUseCase) CreateUser(email, password, firstName, lastName string) (entity.ID, error) {
	s.log.Zap.Info("got to create user")
	e, err := entity.New(email, password, firstName, lastName)
	if err != nil {
		s.log.Zap.Error(err.Error())

		return e.ID, err
	}
	return s.repo.Create(e)
}

// GetUser gets an user
func (s *userUseCase) GetUser(id string) (*entity.User, error) {
	uID, _ := entity.IDFromString(id)
	return s.repo.Get(uID)
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
func (s *userUseCase) DeleteUser(id string) error {
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
	uID, _ := entity.IDFromString(id)
	return s.repo.Delete(uID)
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
