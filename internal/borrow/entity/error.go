package entity

import "errors"

// ErrNotEnoughBooks cannot borrow
var ErrNotEnoughBooks = errors.New("Not enough books")

// ErrBookAlreadyBorrowed cannot borrow
var ErrBookAlreadyBorrowed = errors.New("Book already borrowed")

// ErrBookNotBorrowed cannot return
var ErrBookNotBorrowed = errors.New("Book not borrowed")
