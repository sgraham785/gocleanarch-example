package entity

import (
	"time"

	"github.com/rs/xid"
	bookEntity "github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"golang.org/x/crypto/bcrypt"
)

// ID is id for user
type ID = xid.ID

// NewID create a new book entity ID
func NewID() ID {
	return xid.New()
}

// User entity
type User struct {
	ID        ID
	Email     string
	Password  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Books     []bookEntity.ID
}

// New creates a new user
func New(email, password, firstName, lastName string) (*User, error) {
	u := &User{
		ID:        xid.New(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
	}
	pwd, err := generatePassword(password)
	if err != nil {
		return nil, err
	}
	u.Password = pwd
	err = u.Validate()
	if err != nil {
		return nil, ErrInvalidUserEntity
	}
	return u, nil
}

// AddBook add a book to user
func (u *User) AddBook(id ID) error {
	_, err := u.GetBook(id)
	if err == nil {
		return ErrBookAlreadyBorrowed
	}
	u.Books = append(u.Books, id)
	return nil
}

// RemoveBook remove a book from user
func (u *User) RemoveBook(id ID) error {
	for i, j := range u.Books {
		if j == id {
			u.Books = append(u.Books[:i], u.Books[i+1:]...)
			return nil
		}
	}
	return ErrUserBookNotFound
}

// GetBook get a book for user
func (u *User) GetBook(id ID) (ID, error) {
	for _, v := range u.Books {
		if v == id {
			return id, nil
		}
	}
	return id, ErrUserBookNotFound
}

// Validate validate data
func (u *User) Validate() error {
	if u.Email == "" || u.FirstName == "" || u.LastName == "" || u.Password == "" {
		return ErrInvalidUserEntity
	}

	return nil
}

// ValidatePassword validate user password
func (u *User) ValidatePassword(p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return err
	}
	return nil
}

func generatePassword(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
