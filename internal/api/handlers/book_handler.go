package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/api/models"
	"bazar_book_store/internal/database"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"mime/multipart"
	"net/http"
	"strconv"
)

func RegisterBooksRoutes(r chi.Router) {
	r.Post("/books", AdminOnlyMiddleware(Cfg.createBookHandler))
	r.Get("/books", AuthMiddleware(Cfg.getBooksHandler))
	r.Get("/books_details", AuthMiddleware(Cfg.getBooksDetailsHandler))
	r.Post("/book_favorite", AuthMiddleware(Cfg.updateBookFavoriteHandler))
	r.Get("/favorite_books", AuthMiddleware(Cfg.getFavoriteBooksHandler))

}

func (apiCFG *ApiConfig) createBookHandler(w http.ResponseWriter, r *http.Request, _ database.User) {
	vendorIDStr := r.FormValue("vendor_id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	rateStr := r.FormValue("rate")
	categoriesStr := r.FormValue("categories")
	authorsStr := r.FormValue("authors_ids")

	vendorID, _ := strconv.Atoi(vendorIDStr)
	price, _ := strconv.ParseFloat(priceStr, 64)
	rate, _ := strconv.ParseFloat(rateStr, 64)

	categoriesIDs := helpers.ParseInt32JSON(categoriesStr)
	authorsIDs := helpers.ParseInt32JSON(authorsStr)

	if len(categoriesIDs) == 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "You must add atLeast one category")
		return
	}
	if len(authorsIDs) == 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "You must add atLeast one author")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Image isn't exists")
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	imageURL, err := helpers.UploadImage(file)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not upload the image")
		return
	}

	db := apiCFG.DB

	book, err := db.CreateBook(r.Context(), database.CreateBookParams{
		VendorID:    int32(vendorID),
		Title:       title,
		Description: description,
		Price:       strconv.FormatFloat(price, 'f', 2, 64),
		Rate:        strconv.FormatFloat(rate, 'f', 2, 64),
		AvatarUrl:   imageURL,
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

	for _, categoryID := range categoriesIDs {
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
	for _, authorID := range authorsIDs {
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

func (apiCFG *ApiConfig) getBooksHandler(w http.ResponseWriter, r *http.Request, _ database.User) {

	page, limit, offset := helpers.GetPaginationFromRequest(r)

	db := apiCFG.DB
	totalItems, err := db.CountBooks(r.Context())

	books, err := db.GetBooks(r.Context(), database.GetBooksParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not fetch books")
		return
	}

	response := map[string]interface{}{
		"books":      models.DBBooksToBooks(books),
		"pagination": models.ToPaginated(page, limit, int(totalItems)),
		"message":    "Books gotten successfully",
	}
	helpers.RespondWithJSON(w, http.StatusOK, response)
}

func (apiCFG *ApiConfig) getBooksDetailsHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	categoryID := helpers.StringToNullInt32(r.URL.Query().Get("category_id"))
	vendorID := helpers.StringToNullInt32(r.URL.Query().Get("vendor_id"))
	authorID := helpers.StringToNullInt32(r.URL.Query().Get("author_id"))
	bookID := helpers.StringToNullInt32(r.URL.Query().Get("book_id"))

	db := apiCFG.DB

	booksDetails, err := db.GetBooksDetails(r.Context(), database.GetBooksDetailsParams{
		CategoryID: categoryID,
		VendorID:   vendorID,
		AuthorID:   authorID,
		BookID:     bookID,
		UserID:     user.ID,
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not fetch books")
		return
	}

	books := models.DBBooksDetailsToBooksDetails(booksDetails)

	if len(books) == 1 && bookID.Valid {
		helpers.RespondWithJSON(w, http.StatusOK, &books[0])
	} else {
		if len(books) == 0 {
			books = []models.BookDetails{}
		}
		response := map[string]interface{}{
			"books":   books,
			"message": "Books gotten successfully",
		}
		helpers.RespondWithJSON(w, http.StatusOK, response)
	}
}

func (apiCFG *ApiConfig) updateBookFavoriteHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		BookID int32  `json:"book_id"`
		Action string `json:"action"`
	}

	params, ok := helpers.DecodeBody[parameters](w, r)
	if !ok {
		return
	}

	db := apiCFG.DB
	var err error
	switch params.Action {
	case "add":
		err = db.AddBookFavorite(r.Context(), database.AddBookFavoriteParams{
			UserID: user.ID,
			BookID: params.BookID,
		})
	case "remove":
		err = db.RemoveBookFavorite(r.Context(), database.RemoveBookFavoriteParams{
			UserID: user.ID,
			BookID: params.BookID,
		})
	default:
		helpers.RespondWithError(w, http.StatusInternalServerError, "action must be add or remove.")
		return
	}

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not add/remove book to favorites list")
		return
	}

	response := map[string]interface{}{}

	if params.Action == "add" {
		response["message"] = "Book added to favorites successfully"
		helpers.RespondWithJSON(w, http.StatusCreated, response)
	} else {
		response["message"] = "Book removed from favorites successfully"
		helpers.RespondWithJSON(w, http.StatusOK, response)
	}

}

func (apiCFG *ApiConfig) getFavoriteBooksHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	db := apiCFG.DB

	favBooks, err := db.GetFavoriteBooks(r.Context(), user.ID)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not fetch fav books")
		return
	}

	response := map[string]interface{}{
		"favorite_books": models.DBFavoriteBooksToBooks(favBooks),
		"message":        "Favorite books gotten successfully",
	}
	helpers.RespondWithJSON(w, http.StatusOK, response)
}
