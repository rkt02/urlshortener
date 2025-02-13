package main

import (
	//"fmt"
	//"log"

	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rkt02/urlshortener/internal/db"
	"github.com/rkt02/urlshortener/internal/handlers"
)

func main() {
	connstr := "postgres://postgres:1234@localhost:5432/urlshortdb?sslmode=disable"

	dbInstance := db.OpenPostgresDB(connstr)
	defer dbInstance.Close()

	//testing things that should later be put in handling
	//db.CreateURLMapping(dbInstance, "WSJ.com", utils.EncodeBase62(123545))
	//db.PrintURLTable(dbInstance)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	handler := handlers.NewHandler(dbInstance)

	router.Post("/shorten/{long}", handler.ShortenURL)
	router.Get("/{short}", handler.Redirect)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error Starting Server: %v", err)
	}

}
