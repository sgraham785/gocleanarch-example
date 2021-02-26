package usecase_test

import (
	"testing"
	"time"

	"github.com/sgraham785/gocleanarch-example/internal/book/entity"
	"github.com/sgraham785/gocleanarch-example/internal/book/infrastructure"
	"github.com/sgraham785/gocleanarch-example/internal/book/usecase"
	"github.com/stretchr/testify/assert"
)

func newFixtureBook() *entity.Book {
	return &entity.Book{
		Title:     "I Am Ozzy",
		Author:    "Ozzy Osbourne",
		Pages:     294,
		Quantity:  1,
		CreatedAt: time.Now(),
	}
}

func Test_bookUseCase_CreateBook(t *testing.T) {
	repo := infrastructure.NewInMemRepo()
	m := usecase.New(repo)
	u := newFixtureBook()
	_, err := m.CreateBook(u.Title, u.Author, u.Pages, u.Quantity)
	assert.Nil(t, err)
	assert.False(t, u.CreatedAt.IsZero())
}

func Test_bookUseCase_SearchBooks(t *testing.T) {
	repo := infrastructure.NewInMemRepo()
	m := usecase.New(repo)
	u1 := newFixtureBook()
	u2 := newFixtureBook()
	u2.Title = "Lemmy: Biography"

	uID, _ := m.CreateBook(u1.Title, u1.Author, u1.Pages, u1.Quantity)
	_, _ = m.CreateBook(u2.Title, u2.Author, u2.Pages, u2.Quantity)

	t.Run("search", func(t *testing.T) {
		c, err := m.SearchBooks("ozzy")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(c))
		assert.Equal(t, "I Am Ozzy", c[0].Title)

		c, err = m.SearchBooks("dio")
		assert.Equal(t, entity.ErrBookNotFound, err)
		assert.Nil(t, c)
	})
	t.Run("list all", func(t *testing.T) {
		all, err := m.ListBooks()
		assert.Nil(t, err)
		assert.Equal(t, 2, len(all))
	})

	t.Run("get", func(t *testing.T) {
		saved, err := m.GetBook(uID)
		assert.Nil(t, err)
		assert.Equal(t, u1.Title, saved.Title)
	})
}

func Test_bookUseCase_UpdateBook(t *testing.T) {
	repo := infrastructure.NewInMemRepo()
	m := usecase.New(repo)
	u := newFixtureBook()
	id, err := m.CreateBook(u.Title, u.Author, u.Pages, u.Quantity)
	assert.Nil(t, err)
	saved, _ := m.GetBook(id)
	saved.Title = "Lemmy: Biography"
	assert.Nil(t, m.UpdateBook(saved))
	updated, err := m.GetBook(id)
	assert.Nil(t, err)
	assert.Equal(t, "Lemmy: Biography", updated.Title)
}

func Test_bookUseCase_DeleteBook(t *testing.T) {
	repo := infrastructure.NewInMemRepo()
	m := usecase.New(repo)
	u1 := newFixtureBook()
	u2 := newFixtureBook()
	u2ID, _ := m.CreateBook(u2.Title, u2.Author, u2.Pages, u2.Quantity)

	err := m.DeleteBook(u1.ID)
	assert.Equal(t, entity.ErrBookNotFound, err)

	err = m.DeleteBook(u2ID)
	assert.Nil(t, err)
	_, err = m.GetBook(u2ID)
	assert.Equal(t, entity.ErrBookNotFound, err)
}
