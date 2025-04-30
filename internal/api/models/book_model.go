package models

import (
	"bazar_book_store/internal/database"
	"strconv"
	"time"
)

type Book struct {
	ID          int32     `json:"id"`
	VendorID    int32     `json:"vendor_id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	Rate        float64   `json:"rate,omitempty"`
	IsFavorite  *bool     `json:"is_favorite,omitempty"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type BookDetails struct {
	Book       Book       `json:"book"`
	Authors    []Author   `json:"authors"`
	Categories []Category `json:"categories"`
}

func DBBookToBook(dbBook database.Book) Book {
	rateValue, err := strconv.ParseFloat(dbBook.Rate, 64)
	if err != nil {
		rateValue = 1.0
	}

	priceValue, err := strconv.ParseFloat(dbBook.Price, 64)
	if err != nil {
		priceValue = 0.0
	}

	return Book{
		ID:          dbBook.ID,
		VendorID:    dbBook.VendorID,
		Title:       dbBook.Title,
		Description: dbBook.Description,
		Price:       priceValue,
		Rate:        rateValue,
		CreatedAt:   dbBook.CreatedAt,
		UpdatedAt:   dbBook.UpdatedAt,
		AvatarURL:   dbBook.AvatarUrl,
	}
}

func DBBooksToBooks(dbBooks []database.Book) []Book {
	books := make([]Book, len(dbBooks))

	for i, dbBook := range dbBooks {
		books[i] = DBBookToBook(dbBook)
	}

	return books
}

func DBFavoriteBooksToBooks(dbBooks []database.GetFavoriteBooksRow) []Book {
	favorites := make([]Book, len(dbBooks))

	for i, dbBook := range dbBooks {
		favorites[i] = DBBookToBook(database.Book{
			ID:        dbBook.ID,
			Title:     dbBook.Title,
			Price:     dbBook.Price,
			AvatarUrl: dbBook.AvatarUrl,
		})
	}

	return favorites
}

func DBBooksDetailsToBooksDetails(dbBooksDetails []database.GetBooksDetailsRow) []BookDetails {
	// Create a map to hold unique books by BookID (no duplicates)
	booksMap := make(map[int32]BookDetails)
	// Create maps to group authors and categories by BookID
	bookAuthors := make(map[int32][]Author)
	bookCategories := make(map[int32][]Category)
	bookCategoryIDs := make(map[int32]map[int32]struct{})
	bookAuthorIDs := make(map[int32]map[int32]struct{})

	// Loop through dbBooksDetails and accumulate data
	for _, dbBookDetails := range dbBooksDetails {
		bookID := dbBookDetails.BookID

		// Create the Book object only if it doesn't exist already in booksMap
		if _, exists := booksMap[bookID]; !exists {
			bookDetails := BookDetails{
				Book: DBBookToBook(database.Book{
					ID:          dbBookDetails.BookID,
					VendorID:    dbBookDetails.VendorID,
					Title:       dbBookDetails.Title,
					Description: dbBookDetails.Description,
					Price:       dbBookDetails.Price,
					Rate:        dbBookDetails.Rate,
					AvatarUrl:   dbBookDetails.BookAvatarUrl,
					CreatedAt:   dbBookDetails.CreatedAt,
					UpdatedAt:   dbBookDetails.UpdatedAt,
				}),
			}
			bookDetails.Book.IsFavorite = &dbBookDetails.IsFavorite
			booksMap[bookID] = bookDetails
		}

		// Add author to map of authors for the current book
		if dbBookDetails.AuthorID.Valid {
			authorID := dbBookDetails.AuthorID.Int32
			if _, ok := bookAuthorIDs[bookID]; !ok {
				bookAuthorIDs[bookID] = make(map[int32]struct{})
			}
			if _, exists := bookAuthorIDs[bookID][authorID]; !exists {
				author := DBAuthorToAuthor(database.Author{
					ID:               authorID,
					Name:             dbBookDetails.AuthorName.String,
					ShortDescription: dbBookDetails.AuthorShortDescription.String,
					About:            dbBookDetails.AuthorAbout.String,
					AvatarUrl:        dbBookDetails.AuthorAvatarUrl,
					Rate:             dbBookDetails.AuthorRate.String,
					AuthorType:       dbBookDetails.AuthorType.AuthorTypeEnum,
				})
				bookAuthors[bookID] = append(bookAuthors[bookID], author)
				bookAuthorIDs[bookID][authorID] = struct{}{}
			}

		}

		// Add category to map of categories for the current book

		if dbBookDetails.CategoryID.Valid {
			categoryID := dbBookDetails.CategoryID.Int32
			if _, ok := bookCategoryIDs[bookID]; !ok {
				bookCategoryIDs[bookID] = make(map[int32]struct{})
			}

			if _, exists := bookCategoryIDs[bookID][categoryID]; !exists {
				bookCategories[bookID] = append(bookCategories[bookID], DBCategoryToCategory(database.Category{
					ID:   categoryID,
					Name: dbBookDetails.CategoryName.String,
				}))
				bookCategoryIDs[bookID][categoryID] = struct{}{}
			}
		}
	}

	// Now fill in the authors and categories for each unique book
	var books []BookDetails
	for bookID, bookDetails := range booksMap {
		bookDetails.Authors = bookAuthors[bookID]
		bookDetails.Categories = bookCategories[bookID]
		if len(bookDetails.Authors) == 0 {
			bookDetails.Authors = []Author{}
		}
		if len(bookDetails.Categories) == 0 {
			bookDetails.Categories = []Category{}
		}
		books = append(books, bookDetails)
	}

	return books
}
