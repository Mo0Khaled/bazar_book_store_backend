package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/api/models"
	"bazar_book_store/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"

	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterUserRoutes(r chi.Router) {
	r.Post("/register", Cfg.createUserHandler)
}

func (apiCFG *ApiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	db := apiCFG.DB

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.RespondWithError(w, http.StatusForbidden, "Incorrect credentials")
		return
	}

	user, err := db.CreateUser(r.Context(), database.CreateUserParams{
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     user.ID,
		"expiry_date": time.Now().Add(time.Hour * 24 * 7).Unix(), // صلاحية 24 ساعة
	})
	secretKey := os.Getenv("JWT_SECRET")
	userToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"user":    models.DBUserToUser(user),
		"token":   userToken,
		"message": "User created successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}
