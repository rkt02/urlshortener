package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rkt02/urlshortener/internal/auth"
	"github.com/rkt02/urlshortener/internal/cache"
	"github.com/rkt02/urlshortener/internal/db"
	"github.com/rkt02/urlshortener/internal/handlers"
)

var redisAddress = os.Getenv("REDIS_ADDRESS")
var redisPassword = os.Getenv("REDIS_PASSWORD")

func main() {
	//TODO: This should also either be an environment variable or put into a secure data store
	connstr := "postgres://postgres:1234@localhost:5432/urlshortdb?sslmode=disable"

	/*
		err := godotenv.Load()
		if err != nil {
			log.Printf("Failed to load .env File")
		}
		var redisAddress = os.Getenv("REDIS_ADDRESS")
		var redisPassword = os.Getenv("REDIS_PASSWORD")
	*/

	dbInstance := db.OpenPostgresDB(connstr)
	defer dbInstance.Close()

	redisClient, err := cache.OpenRedisClient(redisAddress, redisPassword, 0)
	if err != nil {
		log.Fatalf("Redis Failed to Start: %v", err)
		return
	}

	//testing things that should later be put in handling
	//db.CreateURLMapping(dbInstance, "WSJ.com", utils.EncodeBase62(123545))
	//db.PrintURLTable(dbInstance)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	handler := handlers.NewHandler(dbInstance, redisClient)

	router.Get("/login", handlers.Login)

	router.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware)
		r.Post("/shorten/{long}", handler.ShortenURL)
		r.Delete("/remove/long/{long}", handler.DeleteLong)
		r.Delete("/remove/short/{short}", handler.DeleteShortCode)
	})

	router.Get("/{short}", handler.Redirect)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error Starting Server: %v", err)
	}

}
