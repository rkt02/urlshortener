package main

import (
	//"fmt"
	"log"

	//"github.com/rkt02/urlshortener/internal/utils"
	"github.com/rkt02/urlshortener/internal/db"
)

func main() {
	connstr := "postgres://postgres:1234@localhost:5432/urlshortdb?sslmode=disable"

	//Need to add a defer db.Close() of the opened db

	dbInstance, err := db.OpenPostgresDB(connstr)

	if err != nil {
		log.Fatal(err)
	}

	defer dbInstance.Close()
}
