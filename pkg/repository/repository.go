package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sgraham785/gocleanarch-example/pkg/config"
)

type Repository struct {
	Pg *sqlx.DB
}

// NewPostgresConn is called to connect to a PostgreSQL database
func NewPostgresConn(c config.PostgresConf) *Repository {
	format := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	dataSource := fmt.Sprintf(format,
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB)
	db, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		fmt.Println(err)
	}
	return &Repository{Pg: db}
}
