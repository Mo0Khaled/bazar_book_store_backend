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
		Categories  []int32 `json:"categories"`
		AuthorsIDs  []int32 `json:"authors_ids"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	if len(params.Categories) == 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "You must add atLeast one category")
		return
	}
	if len(params.AuthorsIDs) == 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "You must add atLeast one author")
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

	var invalidCategories []int32
	var validCategories []database.Category

	for _, categoryID := range params.Categories {
		category, err := db.GetCategoryByID(r.Context(), categoryID)
		if err != nil {
			invalidCategories = append(invalidCategories, categoryID)
		} else {
			validCategories = append(validCategories, category)
			err = db.AddBookCategory(r.Context(), database.AddBookCategoryParams{
				BookID:     book.ID,
				CategoryID: category.ID,
			})

			if err != nil {
				// TODO: i think i should think of undo the created book.
				fmt.Println(err)
			}
		}
	}

	var invalidAuthorsIDs []int32
	for _, authorID := range params.AuthorsIDs {
		err := db.AddBookAuthor(r.Context(), database.AddBookAuthorParams{
			BookID:   book.ID,
			AuthorID: authorID,
		})

		if err != nil {
			fmt.Println(err)
			invalidAuthorsIDs = append(invalidAuthorsIDs, authorID)
		}

	}

	response := map[string]interface{}{
		"book":       models.DBBookToBook(book),
		"categories": models.DBCategoriesToCategories(validCategories),
		"message":    "Book created successfully",
	}
	if len(invalidCategories) > 0 || len(invalidAuthorsIDs) > 0 {
		response["warning"] = fmt.Sprintf("The following categories/authors do not exist: %v %v", invalidCategories, invalidAuthorsIDs)
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}
