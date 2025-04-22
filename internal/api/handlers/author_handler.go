package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/api/models"
	"bazar_book_store/internal/database"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func RegisterAuthorRoutes(r chi.Router) {
	r.Post("/authors", AdminOnlyMiddleware(Cfg.createAuthorHandler))

}

func (apiCFG *ApiConfig) createAuthorHandler(w http.ResponseWriter, r *http.Request, _ database.User) {
	type parameters struct {
		Name             string  `json:"name"`
		ShortDescription string  `json:"short_description"`
		About            string  `json:"about"`
		AuthorType       string  `json:"author_type"`
		AvatarURL        *string `json:"avatar_url"`
		Rate             float64 `json:"rate"`
	}

	params, ok := helpers.DecodeBody[parameters](w, r)
	if !ok {
		return
	}

	db := apiCFG.DB

	author, err := db.CreateAuthor(r.Context(), database.CreateAuthorParams{
		Name:             params.Name,
		ShortDescription: params.ShortDescription,
		About:            params.About,
		AuthorType:       database.AuthorTypeEnum(params.AuthorType),
		AvatarUrl:        helpers.ToNullString(params.AvatarURL),
		Rate:             strconv.FormatFloat(params.Rate, 'f', 2, 64),
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not create author")
		return
	}

	response := map[string]interface{}{
		"author":  models.DBAuthorToAuthor(author),
		"message": "Author created successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}
