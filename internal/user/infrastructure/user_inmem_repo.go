package infrastructure

import (
	"strings"
	"sync"

	"github.com/sgraham785/gocleanarch-example/internal/user/entity"
)

type userInMemRepo struct {
	mtx sync.RWMutex
	m   map[entity.ID]*entity.User
}

// NewInMemRepo create book in memory repository
func NewInMemRepo() UserRepo {
	var m = map[entity.ID]*entity.User{}
	return &userInMemRepo{
		m: m,
	}
}

// Create an user
func (r *userInMemRepo) Create(e *entity.User) (entity.ID, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.m[e.ID] = e
	return e.ID, nil
}

// Get an user
func (r *userInMemRepo) Get(id entity.ID) (*entity.User, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.m[id] == nil {
		return nil, entity.ErrUserNotFound
	}
	return r.m[id], nil
}

// Update an user
func (r *userInMemRepo) Update(e *entity.User) error {
	_, err := r.Get(e.ID)
	if err != nil {
		return err
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.m[e.ID] = e
	return nil
}

// Search users
func (r *userInMemRepo) Search(query string) ([]*entity.User, error) {
	var d []*entity.User
	r.mtx.Lock()
	defer r.mtx.Unlock()
	for _, j := range r.m {
		if strings.Contains(strings.ToLower(j.FirstName), query) {
			d = append(d, j)
		}
	}
	if len(d) == 0 {
		return nil, entity.ErrUserNotFound
	}

	return d, nil
}

// List users
func (r *userInMemRepo) List() ([]*entity.User, error) {
	var d []*entity.User
	r.mtx.Lock()
	defer r.mtx.Unlock()
	for _, j := range r.m {
		d = append(d, j)
	}
	return d, nil
}

// Delete an user
func (r *userInMemRepo) Delete(id entity.ID) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.m[id] == nil {
		return entity.ErrUserNotFound
	}
	r.m[id] = nil
	return nil
}
