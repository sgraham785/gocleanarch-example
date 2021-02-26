package entity_test

import (
	"testing"

	bookEntity "github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/sgraham785/gocleanarch-example/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	u, err := entity.New("sjobs@apple.com", "new_password", "Steve", "Jobs")
	assert.Nil(t, err)
	assert.Equal(t, u.FirstName, "Steve")
	assert.NotNil(t, u.ID)
	assert.NotEqual(t, u.Password, "new_password")
}

func TestUser_ValidatePassword(t *testing.T) {
	u, _ := entity.New("sjobs@apple.com", "new_password", "Steve", "Jobs")
	err := u.ValidatePassword("new_password")
	assert.Nil(t, err)
	err = u.ValidatePassword("wrong_password")
	assert.NotNil(t, err)

}

func TestUser_AddBook(t *testing.T) {
	u, _ := entity.New("sjobs@apple.com", "new_password", "Steve", "Jobs")
	bID := bookEntity.NewID()
	err := u.AddBook(bID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(u.Books))
	err = u.AddBook(bID)
	assert.Equal(t, entity.ErrBookAlreadyBorrowed, err)

}

func TestUser_RemoveBook(t *testing.T) {
	u, _ := entity.New("sjobs@apple.com", "new_password", "Steve", "Jobs")
	err := u.RemoveBook(bookEntity.NewID())
	assert.Equal(t, entity.ErrUserBookNotFound, err)
	bID := bookEntity.NewID()
	_ = u.AddBook(bID)
	err = u.RemoveBook(bID)
	assert.Nil(t, err)
}

func TestUser_GetBook(t *testing.T) {
	u, _ := entity.New("sjobs@apple.com", "new_password", "Steve", "Jobs")
	bID := bookEntity.NewID()
	_ = u.AddBook(bID)
	id, err := u.GetBook(bID)
	assert.Nil(t, err)
	assert.Equal(t, id, bID)
	_, err = u.GetBook(bookEntity.NewID())
	assert.Equal(t, entity.ErrUserBookNotFound, err)
}

func TestUser_Validate(t *testing.T) {
	type test struct {
		email     string
		password  string
		firstName string
		lastName  string
		want      error
	}

	tests := []test{
		{
			email:     "sjobs@apple.com",
			password:  "new_password",
			firstName: "Steve",
			lastName:  "Jobs",
			want:      nil,
		},
		{
			email:     "",
			password:  "new_password",
			firstName: "Steve",
			lastName:  "Jobs",
			want:      entity.ErrInvalidUserEntity,
		},
		{
			email:     "sjobs@apple.com",
			password:  "",
			firstName: "Steve",
			lastName:  "Jobs",
			want:      nil,
		},
		{
			email:     "sjobs@apple.com",
			password:  "new_password",
			firstName: "",
			lastName:  "Jobs",
			want:      entity.ErrInvalidUserEntity,
		},
		{
			email:     "sjobs@apple.com",
			password:  "new_password",
			firstName: "Steve",
			lastName:  "",
			want:      entity.ErrInvalidUserEntity,
		},
	}
	for _, tc := range tests {

		_, err := entity.New(tc.email, tc.password, tc.firstName, tc.lastName)
		assert.Equal(t, err, tc.want)
	}

}
