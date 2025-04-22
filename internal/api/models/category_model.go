package models

import (
	"bazar_book_store/internal/database"
	"time"
)

type Category struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DBCategoryToCategory(dbCategory database.Category) Category {
	return Category{
		ID:        dbCategory.ID,
		Name:      dbCategory.Name,
		CreatedAt: dbCategory.CreatedAt,
		UpdatedAt: dbCategory.UpdatedAt,
	}
}
