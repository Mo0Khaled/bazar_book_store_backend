package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/api/models"
	"bazar_book_store/internal/database"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"net/http"
	"strconv"
)

func RegisterBooksRoutes(r chi.Router) {
	r.Post("/books", AdminOnlyMiddleware(Cfg.createBookHandler))

}

func (apiCFG *ApiConfig) createBookHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		VendorID    int32   `json:"vendor_id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Rate        float64 `json:"rate"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	db := apiCFG.DB

	book, err := db.CreateBook(r.Context(), database.CreateBookParams{
		VendorID:    params.VendorID,
		Title:       params.Title,
		Description: params.Description,
		Price:       strconv.FormatFloat(params.Price, 'f', 2, 64),
		Rate:        strconv.FormatFloat(params.Rate, 'f', 2, 64),
	})

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			fmt.Println(pqErr.Code.Name())
			fmt.Println(pqErr.Constraint)
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				if pqErr.Constraint == "books_vendor_id_fkey" {
					helpers.RespondWithError(w, http.StatusBadRequest, "Vendor ID isn't exists")
					return
				}
			}
		}
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not create book")
		return
	}

	response := map[string]interface{}{
		"book":    models.DBBookToBook(book),
		"message": "Book created successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}
