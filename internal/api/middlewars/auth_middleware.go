package middlewars

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/api/handlers"
	"bazar_book_store/internal/database"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

type contextKey string

const userIDKey contextKey = "user_id"

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func AuthMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid claims")
			return
		}

		userIDFloat, ok := claims["user_id"].(int32)
		if !ok {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		user, err := handlers.Cfg.DB.GetUser(r.Context(), userIDFloat)
		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID")
			return
		}
		handler(w, r, user)
	}
}
