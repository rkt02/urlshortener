package main

import (
	//"fmt"
	//"log"

	//"github.com/rkt02/urlshortener/internal/utils"
	"github.com/rkt02/urlshortener/internal/db"
)

func main() {
	connstr := "postgres://postgres:1234@localhost:5432/urlshortdb?sslmode=disable"

	dbInstance := db.OpenPostgresDB(connstr)
	defer dbInstance.Close()

	//testing things that should later be put in handling
	db.CreateURLMapping(dbInstance, "gogkkejj", "aaa")
	db.PrintURLTable(dbInstance)

}
