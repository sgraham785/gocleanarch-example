package infrastructure

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sgraham785/gocleanarch-example/internal/user/entity"
)

type userPgRepo struct {
	db *sqlx.DB
}

// NewPgRepo create new user repository
func NewPgRepo(db *sqlx.DB) UserRepo {
	return &userPgRepo{
		db: db,
	}
}

// Create an user
func (r *userPgRepo) Create(e *entity.User) (entity.ID, error) {
	sql := `insert into "user" (id, email, password, first_name, last_name, created_at) values($1,$2,$3,$4,$5,$6)`
	fmt.Printf("Create Repo: user=%v\n", e)
	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.Email,
		e.Password,
		e.FirstName,
		e.LastName,
		time.Now().Format("2006-01-02"),
	)
	if err != nil {
		return e.ID, err
	}
	err = stmt.Close()
	if err != nil {
		return e.ID, err
	}
	return e.ID, nil
}

// Get an user
func (r *userPgRepo) Get(id entity.ID) (*entity.User, error) {
	sql := `select id, email, first_name, last_name, created_at from "user" where id = $1`
	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	var u entity.User
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)
	}
	stmt, err = r.db.Prepare(`select book_id from book_user where user_id = ?`)
	if err != nil {
		return nil, err
	}
	rows, err = stmt.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i entity.ID
		err = rows.Scan(&i)
		u.Books = append(u.Books, i)
	}
	return &u, nil
}

// Update an user
func (r *userPgRepo) Update(e *entity.User) error {
	sql := `update "user" set email = $1, password = $2, first_name = $3, last_name = $4, updated_at = $5 where id = $6`

	e.UpdatedAt = time.Now()
	_, err := r.db.Exec(sql, e.Email, e.Password, e.FirstName, e.LastName, e.UpdatedAt.Format("2006-01-02"), e.ID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("delete from book_user where user_id = ?", e.ID)
	if err != nil {
		return err
	}
	for _, b := range e.Books {
		_, err := r.db.Exec("insert into book_user values(?,?,?)", e.ID, b, time.Now().Format("2006-01-02"))
		if err != nil {
			return err
		}
	}
	return nil
}

// Search users
func (r *userPgRepo) Search(query string) ([]*entity.User, error) {
	sql := `select id from "user" where name like $1`

	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var ids []entity.ID
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i entity.ID
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("not found")
	}
	var users []*entity.User
	for _, id := range ids {
		u, err := r.Get(id)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// List users
func (r *userPgRepo) List() ([]*entity.User, error) {
	sql := `select id from "user"`
	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var ids []entity.ID
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i entity.ID
		err = rows.Scan(&i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("not found")
	}
	var users []*entity.User
	for _, id := range ids {
		u, err := r.Get(id)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Delete an user
func (r *userPgRepo) Delete(id entity.ID) error {
	sql := `delete from "user" where id = $1`
	_, err := r.db.Exec(sql, id)
	if err != nil {
		return err
	}
	return nil
}
