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
