package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rkt02/urlshortener/internal/db"
	"github.com/rkt02/urlshortener/internal/utils"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(dbconn *sql.DB) *Handler {
	return &Handler{DB: dbconn}
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

	long_url, err := db.GetLongFromShort(h.DB, short_code)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("URL Not Found: %v", err)
			http.Error(w, "URL not found", http.StatusNotFound)
		} else {
			log.Printf("DB Retrieve Long From ShortError: %v", err)
			http.Error(w, "Database Error", http.StatusInternalServerError)
		}
		return
	}

	if !strings.HasPrefix(long_url, "http://") && !strings.HasPrefix(long_url, "https://") {
		// If not, prepend a default scheme. You can customize this logic.
		long_url = "http://" + long_url
	}

	http.Redirect(w, r, long_url, http.StatusFound)
}
