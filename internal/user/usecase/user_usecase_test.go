package usecase_test

import (
	"testing"
	"time"

	"github.com/sgraham785/gocleanarch-example/internal/user/entity"
	"github.com/sgraham785/gocleanarch-example/internal/user/infrastructure"
	"github.com/sgraham785/gocleanarch-example/internal/user/usecase"
	"github.com/sgraham785/gocleanarch-example/pkg/logger"
	"github.com/sgraham785/gocleanarch-example/pkg/server"
	"github.com/stretchr/testify/assert"
)

func newFixtureUser() *entity.User {
	return &entity.User{
		ID:        entity.NewID(),
		Email:     "ozzy@metalgods.net",
		Password:  "123456",
		FirstName: "Ozzy",
		LastName:  "Osbourne",
		CreatedAt: time.Now(),
	}
}

func Test_userUseCase_CreateUser(t *testing.T) {
	r := infrastructure.NewInMemRepo()
	logger := logger.New()
	defer logger.Zap.Sync()
	s := &server.Server{
		Log: logger,
	}
	uc := usecase.New(s, r)
	u := newFixtureUser()
	_, err := uc.CreateUser(u.Email, u.Password, u.FirstName, u.LastName)
	assert.Nil(t, err)
	assert.False(t, u.CreatedAt.IsZero())
	assert.True(t, u.UpdatedAt.IsZero())
}

func Test_userUseCase_SearchUsers(t *testing.T) {
	r := infrastructure.NewInMemRepo()
	logger := logger.New()
	defer logger.Zap.Sync()
	s := &server.Server{
		Log: logger,
	}
	uc := usecase.New(s, r)
	u1 := newFixtureUser()
	u2 := newFixtureUser()
	u2.FirstName = "Lemmy"

	uID, _ := uc.CreateUser(u1.Email, u1.Password, u1.FirstName, u1.LastName)
	_, _ = uc.CreateUser(u2.Email, u2.Password, u2.FirstName, u2.LastName)

	t.Run("search", func(t *testing.T) {
		c, err := uc.SearchUsers("ozzy")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(c))
		assert.Equal(t, "Osbourne", c[0].LastName)

		c, err = uc.SearchUsers("dio")
		assert.Equal(t, entity.ErrUserNotFound, err)
		assert.Nil(t, c)
	})
	t.Run("list all", func(t *testing.T) {
		all, err := uc.ListUsers()
		assert.Nil(t, err)
		assert.Equal(t, 2, len(all))
	})

	t.Run("get", func(t *testing.T) {
		saved, err := uc.GetUser(uID.String())
		assert.Nil(t, err)
		assert.Equal(t, u1.FirstName, saved.FirstName)
	})
}

func Test_userUseCase_UpdateUser(t *testing.T) {
	r := infrastructure.NewInMemRepo()
	logger := logger.New()
	defer logger.Zap.Sync()
	s := &server.Server{
		Log: logger,
	}
	uc := usecase.New(s, r)
	u := newFixtureUser()
	id, err := uc.CreateUser(u.Email, u.Password, u.FirstName, u.LastName)
	assert.Nil(t, err)
	saved, _ := uc.GetUser(id.String())
	saved.FirstName = "Dio"
	saved.Books = append(saved.Books, entity.NewID())
	assert.Nil(t, uc.UpdateUser(saved))
	updated, err := uc.GetUser(id.String())
	assert.Nil(t, err)
	assert.Equal(t, "Dio", updated.FirstName)
	assert.False(t, updated.UpdatedAt.IsZero())
	assert.Equal(t, 1, len(updated.Books))
}

func Test_userUseCase_DeleteUser(t *testing.T) {
	r := infrastructure.NewInMemRepo()
	logger := logger.New()
	defer logger.Zap.Sync()
	s := &server.Server{
		Log: logger,
	}
	uc := usecase.New(s, r)
	u1 := newFixtureUser()
	u2 := newFixtureUser()
	u2ID, _ := uc.CreateUser(u2.Email, u2.Password, u2.FirstName, u2.LastName)

	err := uc.DeleteUser(u1.ID.String())
	assert.Equal(t, entity.ErrUserNotFound, err)

	err = uc.DeleteUser(u2ID.String())
	assert.Nil(t, err)
	_, err = uc.GetUser(u2ID.String())
	assert.Equal(t, entity.ErrUserNotFound, err)

	u3 := newFixtureUser()
	id, _ := uc.CreateUser(u3.Email, u3.Password, u3.FirstName, u3.LastName)
	saved, _ := uc.GetUser(id.String())
	saved.Books = []entity.ID{entity.NewID()}
	_ = uc.UpdateUser(saved)
	err = uc.DeleteUser(id.String())
	assert.Equal(t, entity.ErrUserCannotBeDeleted, err)
}
