package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/api/models"
	"bazar_book_store/internal/database"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func RegisterCategoryRoutes(r chi.Router, h *Handler) {
	r.Post("/categories", AdminOnlyMiddleware(h.createCategoryHandler))
	r.Get("/category/{categoryID}", AuthMiddleware(h.getCategoryHandler))

}

func (h *Handler) createCategoryHandler(w http.ResponseWriter, r *http.Request, _ database.User) {
	type parameters struct {
		Name string `json:"name"`
	}

	params, ok := helpers.DecodeBody[parameters](w, r)
	if !ok {
		return
	}

	db := h.Cfg.DB

	category, err := db.CreateCategory(r.Context(), params.Name)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not create category")
		return
	}

	response := map[string]interface{}{
		"category": models.DBCategoryToCategory(category),
		"message":  "Book created successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}

func (h *Handler) getCategoryHandler(w http.ResponseWriter, r *http.Request, _ database.User) {
	categoryIDStr := chi.URLParam(r, "categoryID")

	categoryID, err := strconv.Atoi(categoryIDStr)

	if err != nil {
		helpers.RespondWithError(w, http.StatusForbidden, "wrong category ID")
		return
	}

	db := h.Cfg.DB

	category, err := db.GetCategoryByID(r.Context(), int32(categoryID))

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not fetch category")
		return
	}

	response := map[string]interface{}{
		"category": models.DBCategoryToCategory(category),
		"message":  "Category gotten successfully",
	}
	helpers.RespondWithJSON(w, http.StatusOK, response)
}
