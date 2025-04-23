package models

import (
	"bazar_book_store/internal/database"
	"strconv"
	"time"
)

type Book struct {
	ID          int32     `json:"id"`
	VendorID    int32     `json:"vendor_id"`
	Title       string    `json:"title"`
	Description string    `json:"Description"`
	Price       float64   `json:"price"`
	Rate        float64   `json:"rate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	}
}

func DBBooksToBooks(dbBooks []database.Book) []Book {
	books := make([]Book, len(dbBooks))

	for i, dbBook := range dbBooks {
		books[i] = DBBookToBook(dbBook)
	}

	return books
}

func DBBooksDetailsToBooksDetails(dbBooksDetails []database.GetBooksDetailsRow) []BookDetails {
	// Create a map to hold unique books by BookID (no duplicates)
	booksMap := make(map[int32]BookDetails)
	// Create maps to group authors and categories by BookID
	bookAuthors := make(map[int32][]Author)
	bookCategories := make(map[int32][]Category)

	// Loop through dbBooksDetails and accumulate data
	for _, dbBookDetails := range dbBooksDetails {
		bookID := dbBookDetails.BookID

		// Create the Book object only if it doesn't exist already in booksMap
		if _, exists := booksMap[bookID]; !exists {
			booksMap[bookID] = BookDetails{
				Book: DBBookToBook(database.Book{
					ID:          dbBookDetails.BookID,
					VendorID:    dbBookDetails.VendorID,
					Title:       dbBookDetails.Title,
					Description: dbBookDetails.Description,
					Price:       dbBookDetails.Price,
					Rate:        dbBookDetails.Rate,
					CreatedAt:   dbBookDetails.CreatedAt,
					UpdatedAt:   dbBookDetails.UpdatedAt,
				}),
			}
		}

		// Add author to map of authors for the current book
		if dbBookDetails.AuthorID.Valid {
			author := DBAuthorToAuthor(database.Author{
				ID:               dbBookDetails.AuthorID.Int32,
				Name:             dbBookDetails.AuthorName.String,
				ShortDescription: dbBookDetails.AuthorShortDescription.String,
				About:            dbBookDetails.AuthorAbout.String,
				AvatarUrl:        dbBookDetails.AuthorAvatarUrl,
				Rate:             dbBookDetails.AuthorRate.String,
				AuthorType:       dbBookDetails.AuthorType.AuthorTypeEnum,
			})
			bookAuthors[bookID] = append(bookAuthors[bookID], author)
		}

		// Add category to map of categories for the current book

		if dbBookDetails.CategoryID.Valid {
			category := DBCategoryToCategory(database.Category{
				ID:   dbBookDetails.CategoryID.Int32,
				Name: dbBookDetails.CategoryName.String,
			})
			bookCategories[bookID] = append(bookCategories[bookID], category)
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
