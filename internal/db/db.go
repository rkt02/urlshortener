//should I do error checking within these functions or put the responsibility on main

package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func OpenPostgresDB(connstr string) *sql.DB {
	db, err := sql.Open("postgres", connstr)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	CreateURLTable(db)

	return db
}

// think about whether this should be done in main or with the openpostgresdb function
func CreateURLTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS url(
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(50) UNIQUE,
		long_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	)`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalf("Inside Create URL table func %v", err)
	}
}

// TODO: Add back in error here for caller to be able to respong
// Big case is that thre has already been a short code key added
// May actually not be a problem because short code is produced after the long url is added and a serial
// Id has been assigned TBD
func CreateURLMapping(db *sql.DB, short_code, long_url string) int64 {
	var id int64
	query := `INSERT INTO url (short_code, long_url) VALUES ($1, $2) RETURNING id`

	err := db.QueryRow(query, short_code, long_url).Scan(&id)

	if err != nil {
		log.Fatalf("Inside URL Map Func%v", err)
	}

	return id
}

func PrintURLTable(db *sql.DB) {
	query := "SELECT short_code, long_url FROM url"
	//data := []Url{}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Print URL Select Query%v", err)
	}

	defer rows.Close()

	var short_code string
	var long_url string

	for rows.Next() {
		err := rows.Scan(&short_code, &long_url)
		if err != nil {
			log.Fatalf("Scanning Retrieved rows%v", err)
		}

		fmt.Println(short_code)
		//data = append(data, Url{short_code, long_url})
	}
}
