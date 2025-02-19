package auth

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Missing Authorization Header")
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Printf("Malformed Token")
			http.Error(w, "Malformed token", http.StatusUnauthorized)
			return
		}

		parser := jwt.Parser{SkipClaimsValidation: true}
		token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			log.Printf("Invalid token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("Invalid token claims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			log.Printf("Invalid token exp claim")
			http.Error(w, "Inavlid token exp claim", http.StatusUnauthorized)
			return
		}
		if int64(exp) < time.Now().Unix() {
			log.Printf("Expired Token")
			http.Error(w, "Expired Token", http.StatusUnauthorized)
			return
		}

		if admin, exists := claims["admin"]; !exists || admin == false {
			log.Printf("Insufficient Privelages")
			http.Error(w, "Insufficient Privelages", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)

	})
}
