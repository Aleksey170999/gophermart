package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	usersTable  = "users"
	ordersTable = "orders"
)

func NewPostgresDB(url string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
