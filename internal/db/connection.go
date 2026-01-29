package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgres() (*sql.DB, error) {
	db, err := sql.Open("postgres",
		"postgres://useradmin:userpassword@localhost:5432/userdb?sslmode=disable",
	)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
