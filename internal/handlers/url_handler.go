package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/rkt02/urlshortener/internal/cache"
	"github.com/rkt02/urlshortener/internal/db"
	"github.com/rkt02/urlshortener/internal/utils"
)

type Handler struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewHandler(dbconn *sql.DB, redis *redis.Client) *Handler {
	return &Handler{DB: dbconn, Redis: redis}
}

/*
type shortenRequest struct {
	LongURL string `json:"long_url"`
}
*/

type shortenResponse struct {
	ShortCode string `json:"short_code"`
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	long_url := chi.URLParam(r, "long")

	if long_url == "" {
		http.Error(w, "Empty URL to Shorten", http.StatusBadRequest)
		return
	}

	id, err := db.CreateURLMapping(h.DB, "", long_url)
	if err != nil {
		log.Printf("DB insert error: %v", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	encoded := utils.EncodeBase62(id)

	//find latest insert and change the short code from empty to the encoded string
	err = db.UpdateShortCodeByID(h.DB, id, encoded)
	if err != nil {
		log.Printf("DB Update error: %v", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	keyStr := "url:" + encoded
	//Would change to longer ttl in production, 1 minute here to be able to easily test
	err = cache.SetCache(h.Redis, keyStr, long_url, 1*time.Minute)
	if err != nil {
		log.Printf("Error setting Cache in ShortenURL: %v", err)
	}

	res := shortenResponse{ShortCode: "http://localhost:8080/" + encoded}

	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("Failed to marshal: %v", err)
		http.Error(w, "Marshal Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(resJSON)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	short_code := chi.URLParam(r, "short")

	if short_code == "" {
		http.Error(w, "Empty Short Code for Redirect", http.StatusBadRequest)
		return
	}

	keyStr := "url:" + short_code
	long_url, err := cache.GetCache(h.Redis, keyStr)
	if err != nil {
		if err == redis.Nil {
			long_url, err = db.GetLongFromShort(h.DB, short_code)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Printf("URL Not Found: %v", err)
					http.Error(w, "URL not found", http.StatusNotFound)
					return
				} else {
					log.Printf("DB Retrieve Long From ShortError: %v", err)
					http.Error(w, "Database Retrieve Error", http.StatusInternalServerError)
					return
				}
			}

			err = cache.SetCache(h.Redis, keyStr, long_url, 5*time.Minute)
			if err != nil {
				log.Printf("Error Setting Cache in Redirect: %v", err)
			}
		} else {
			log.Printf("Error Getting cache in Redirect: %v", err)
			http.Error(w, "Cache Retrieve Error", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Cache Hit in Redirect!")
	}

	if !strings.HasPrefix(long_url, "http://") && !strings.HasPrefix(long_url, "https://") {
		// If not, prepend a default scheme. You can customize this logic.
		long_url = "http://" + long_url
	}

	http.Redirect(w, r, long_url, http.StatusFound)
}

func (h *Handler) DeleteLong(w http.ResponseWriter, r *http.Request) {
	long_url := chi.URLParam(r, "long")

	if long_url == "" {
		http.Error(w, "Empty Long Url for Delete", http.StatusBadRequest)
		return
	}

	rowsDeleted, err := db.DeleteAllLong(h.DB, long_url)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Long URL Not Found: %v", err)
			http.Error(w, "Long URL not found", http.StatusNotFound)
		} else {
			log.Printf("DB Delete by Long URL Error: %v", err)
			http.Error(w, "DB Delete by Long URL Error", http.StatusInternalServerError)
		}
		return
	}

	//IF 0 rows deleted have a different response, no rows deleted does not give an error (the desired state is achieved)
	//UPPDATE: To maintain idempotency, all delete operations produce the same end state, even if that means 0 rows were affected

	log.Printf("%d Rows Deleted with URL: %s", rowsDeleted, long_url)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteShortCode(w http.ResponseWriter, r *http.Request) {
	short_code := chi.URLParam(r, "short")

	if short_code == "" {
		http.Error(w, "Empty Long Url for Delete", http.StatusBadRequest)
		return
	}

	rowsDeleted, err := db.DeleteShort(h.DB, short_code)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Long URL Not Found: %v", err)
			http.Error(w, "Long URL not found", http.StatusNotFound)
		} else {
			log.Printf("DB Delete by Long URL Error: %v", err)
			http.Error(w, "DB Delete by Long URL Error", http.StatusInternalServerError)
		}
		return
	}

	keyString := "url:" + short_code
	cacheDeleted, err := cache.DeleteCache(*h.Redis, keyString)

	//IF 0 rows deleted have a different response, no rows deleted does not give an error (the desired state is achieved)
	//UPPDATE: To maintain idempotency, all delete operations produce the same end state, even if that means 0 rows were affected

	log.Printf("%d BD Rows Deleted with Short Code: %s", rowsDeleted, keyString)
	log.Printf("%d Cache Rows Deleted with Short Code: %s", cacheDeleted, keyString)
	w.WriteHeader(http.StatusNoContent)
}
