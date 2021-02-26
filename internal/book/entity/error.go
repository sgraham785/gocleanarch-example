package entity

import "errors"

// ErrBookNotFound not found
var ErrBookNotFound = errors.New("Book not found")

// ErrInvalidBookEntity invalid book entity
var ErrInvalidBookEntity = errors.New("Invalid book entity")

// ErrBookCannotBeDeleted cannot be deleted
var ErrBookCannotBeDeleted = errors.New("Book cannot be deleted")
