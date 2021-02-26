package infrastructure

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sgraham785/gocleanarch-example/internal/book/entity"
)

// bookPgRepo pg database repo
type bookPgRepo struct {
	db *sqlx.DB
}

// NewPgRepo create new book postgres repo
func NewPgRepo(db *sqlx.DB) BookRepo {
	return &bookPgRepo{
		db: db,
	}
}

// Create a book
func (r *bookPgRepo) Create(e *entity.Book) (entity.ID, error) {
	sql := `insert into book (id, title, author, pages, quantity, created_at) 
	values($1,$2,$3,$4,$5,$6)`

	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return e.ID, err
	}
	_, err = stmt.Exec(
		e.ID,
		e.Title,
		e.Author,
		e.Pages,
		e.Quantity,
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

// Get a book
func (r *bookPgRepo) Get(id entity.ID) (*entity.Book, error) {
	sql := `select id, title, author, pages, quantity, created_at from book where id = $1`
	var book entity.Book
	err := r.db.QueryRowx(sql, id).StructScan(&book)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// Update a book
func (r *bookPgRepo) Update(e *entity.Book) error {
	sql := `update book set title = $1, author = $2, pages = $3, quantity = $4, updated_at = $5 where id = $6`
	e.UpdatedAt = time.Now()
	err := r.db.QueryRowx(sql, e.Title, e.Author, e.Pages, e.Quantity, e.UpdatedAt.Format("2006-01-02"), e.ID).Err()
	if err != nil {
		return err
	}
	return nil
}

// Search books
func (r *bookPgRepo) Search(query string) ([]*entity.Book, error) {
	sql := `select id, title, author, pages, quantity, created_at from book where title like $1`
	var books []*entity.Book
	err := r.db.QueryRowx(sql, query).StructScan(&books)
	if err != nil {
		return nil, err
	}
	// for _, b := range books {
	// 	books = append(books, b)
	// }

	return books, nil
}

// List books
func (r *bookPgRepo) List() ([]*entity.Book, error) {
	sql := `select id, title, author, pages, quantity, created_at from book`
	var books []*entity.Book
	err := r.db.QueryRowx(sql).StructScan(books)
	if err != nil {
		return nil, err
	}
	// for _, b := range books {
	// 	books = append(books, b)
	// }

	return books, nil
}

// Delete a book
func (r *bookPgRepo) Delete(id entity.ID) error {
	sql := `delete from book where id = $1`
	_, err := r.db.Exec(sql, id)
	if err != nil {
		return err
	}
	return nil
}
