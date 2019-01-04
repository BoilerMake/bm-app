package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Open starts the connection to the database.
func Open(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, err
}
