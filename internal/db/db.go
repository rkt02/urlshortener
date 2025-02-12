package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func OpenPostgresDB(connstr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connstr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateURLTable(db sql.DB) {

}
