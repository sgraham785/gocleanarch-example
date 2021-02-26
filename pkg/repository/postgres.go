package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sgraham785/gocleanarch-example/config"
)

// Conn cares the repo connections
// type Conn struct {
// 	Pg *sqlx.DB
// }

// NewPostgresConn is called to connect to a PostgreSQL database
func NewPostgresConn(c config.PostgresConf) *sqlx.DB {
	fmt.Printf("PostgresHost=%s\n", c.PostgresHost)
	fmt.Printf("PostgresPort=%d\n", c.PostgresPort)
	format := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	dataSource := fmt.Sprintf(format,
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB)
	db, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		fmt.Println(err)
	}
	// return &Conn{
	// 	Pg: db,
	// }
	return db
}
