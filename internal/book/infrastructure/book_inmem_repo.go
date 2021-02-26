package infrastructure

import (
	"strings"

	"github.com/sgraham785/gocleanarch-example/internal/book/entity"
)

type bookInMemRepo struct {
	m map[entity.ID]*entity.Book
}

// NewInMemRepo create book in memory repository
func NewInMemRepo() BookRepo {
	var m = map[entity.ID]*entity.Book{}
	return &bookInMemRepo{
		m: m,
	}
}

// Create a book
func (r *bookInMemRepo) Create(e *entity.Book) (entity.ID, error) {
	r.m[e.ID] = e
	return e.ID, nil
}

//Get a book
func (r *bookInMemRepo) Get(id entity.ID) (*entity.Book, error) {
	if r.m[id] == nil {
		return nil, entity.ErrBookNotFound
	}
	return r.m[id], nil
}

//Update a book
func (r *bookInMemRepo) Update(e *entity.Book) error {
	_, err := r.Get(e.ID)
	if err != nil {
		return err
	}
	r.m[e.ID] = e
	return nil
}

//Search books
func (r *bookInMemRepo) Search(query string) ([]*entity.Book, error) {
	var d []*entity.Book
	for _, j := range r.m {
		if strings.Contains(strings.ToLower(j.Title), query) {
			d = append(d, j)
		}
	}
	return d, nil
}

//List books
func (r *bookInMemRepo) List() ([]*entity.Book, error) {
	var d []*entity.Book
	for _, j := range r.m {
		d = append(d, j)
	}
	return d, nil
}

//Delete a book
func (r *bookInMemRepo) Delete(id entity.ID) error {
	if r.m[id] == nil {
		return entity.ErrBookNotFound
	}
	r.m[id] = nil
	return nil
}
