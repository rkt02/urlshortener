package handlers

import (
	"database/sql"
	"net/http"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(dbconn *sql.DB) *Handler {
	return &Handler{DB: dbconn}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {

}
