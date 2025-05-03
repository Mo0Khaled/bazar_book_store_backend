package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/database"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func RegisterApiTokensRoutes(r chi.Router) {
	r.Post("/generate-api-token", AdminOnlyMiddleware(Cfg.createApiTokenHandler))

}

func (apiCFG *ApiConfig) createApiTokenHandler(w http.ResponseWriter, r *http.Request, _ database.User) {
	db := apiCFG.DB

	token, err := helpers.GenerateAPIToken()

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Something went wrong!")
		return
	}

	tokenData, err := db.CreateApiToken(r.Context(), database.CreateApiTokenParams{
		ApiToken:     token,
		RequestLimit: 5000,
		ExpiresAt:    time.Now().AddDate(1, 0, 0),
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not create api token")
		return
	}

	response := map[string]interface{}{
		"token_data": tokenData,
		"message":    "Token created successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}
