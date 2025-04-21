package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/database"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

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

		userIDFloat, ok := claims["user_id"].(float64)

		if !ok {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID")
			return
		}
		userID := int32(userIDFloat)

		user, err := Cfg.DB.GetUser(r.Context(), userID)
		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		handler(w, r, user)
	}
}
