package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/rkt02/urlshortener/internal/auth"
)

var adminUsername = os.Getenv("USERNAME")
var adminPassword = os.Getenv("PASSWORD")

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Environment: Username %s, Password %s \n", adminUsername, adminPassword)
	log.Printf("Request: Username %s, Password %s \n", req.Username, req.Password)

	if req.Username != adminUsername || req.Password != adminPassword {
		log.Printf("Invalid Login")
		http.Error(w, "Invalid Login", http.StatusUnauthorized)
		return
	}

	//Make JWT Token
	tokenString, err := auth.GenerateJWT(req.Username)
	if err != nil {
		log.Printf("Failed to Generate JWT: %v", err)
		http.Error(w, "Error Generating Token String", http.StatusInternalServerError)
	}

	res := LoginResponse{Token: tokenString}

	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("Failed to marshal: %v", err)
		http.Error(w, "Marshal Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
	//Return it in the response
}
