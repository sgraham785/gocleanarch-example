package entity

import (
	"time"

	"github.com/rs/xid"
)

// ID is id for book
type ID = xid.ID

// NewID create a new book entity ID
func NewID() ID {
	return xid.New()
}

// Book entity
type Book struct {
	ID        ID
	Title     string
	Author    string
	Pages     int
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// New creates a new book
func New(title string, author string, pages int, quantity int) (*Book, error) {
	b := &Book{
		ID:        xid.New(),
		Title:     title,
		Author:    author,
		Pages:     pages,
		Quantity:  quantity,
		CreatedAt: time.Now(),
	}
	err := b.Validate()
	if err != nil {
		return nil, ErrInvalidBookEntity
	}
	return b, nil
}

// Validate validate book
func (b *Book) Validate() error {
	if b.Title == "" || b.Author == "" || b.Pages <= 0 || b.Quantity <= 0 {
		return ErrInvalidBookEntity
	}
	return nil
}
