package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/database"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
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

		if CheckRateLimit(w, r, user.IsAdmin) == false {
			return
		}

		handler(w, r, user)
	}
}

func AdminOnlyMiddleware(handler authedHandler) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request, user database.User) {
		if !user.IsAdmin {
			helpers.RespondWithError(w, http.StatusForbidden, "You don't have access to this API")
			return
		}
		handler(w, r, user)
	})
}

func CheckRateLimit(w http.ResponseWriter, r *http.Request, isAdmin bool) bool {
	if isAdmin {
		return true
	}

	rdb := Cfg.RDB
	db := Cfg.DB

	token := r.Header.Get("API-Token")
	if token == "" {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Missing API token")
		return false
	}

	apiToken, err := db.GetApiToken(r.Context(), token)

	if err != nil || apiToken.ExpiresAt.Before(time.Now()) {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return false
	}

	monthKey := time.Now().Format("2006-01")
	key := fmt.Sprintf("rate_limit:%s:%s", token, monthKey)

	count, err := rdb.Incr(r.Context(), key).Result()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Rate limit error")
		return false
	}

	if count == 1 {
		now := time.Now()
		firstOfNextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
		ttl := time.Until(firstOfNextMonth)
		rdb.Expire(r.Context(), key, ttl)
	}

	if int32(count) > apiToken.RequestLimit {
		helpers.RespondWithError(w, http.StatusTooManyRequests, "Rate limit exceeded")
		return false
	}
	return true
}
