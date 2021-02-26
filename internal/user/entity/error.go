package entity

import "errors"

// ErrUserNotFound not found
var ErrUserNotFound = errors.New("User not found")

// ErrUserBookNotFound book not found
var ErrUserBookNotFound = errors.New("Book not found for user")

// ErrInvalidUserEntity invalid user entity
var ErrInvalidUserEntity = errors.New("Invalid user entity")

// ErrUserCannotBeDeleted cannot be deleted
var ErrUserCannotBeDeleted = errors.New("User cannot be deleted")

// ErrNotEnoughBooks cannot borrow
var ErrNotEnoughBooks = errors.New("Not enough books")

// ErrBookAlreadyBorrowed cannot borrow
var ErrBookAlreadyBorrowed = errors.New("Book already borrowed")

// ErrBookNotBorrowed cannot return
var ErrBookNotBorrowed = errors.New("Book not borrowed")
