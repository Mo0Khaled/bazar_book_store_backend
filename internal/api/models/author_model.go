package models

import (
	"bazar_book_store/internal/database"
	"strconv"
	"time"
)

type Author struct {
	ID               int32     `json:"id"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	About            string    `json:"about"`
	AuthorType       string    `json:"author_type"`
	AvatarURL        *string   `json:"avatar_url"`
	Rate             float64   `json:"rate"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func DBAuthorToAuthor(dbAuthor database.Author) Author {
	rateValue, err := strconv.ParseFloat(dbAuthor.Rate, 64)
	if err != nil {
		rateValue = 1.0
	}
	var avatarURL *string
	if dbAuthor.AvatarUrl.Valid {
		avatarURL = &dbAuthor.AvatarUrl.String
	}
	return Author{
		ID:               dbAuthor.ID,
		Name:             dbAuthor.Name,
		ShortDescription: dbAuthor.ShortDescription,
		About:            dbAuthor.About,
		AuthorType:       string(dbAuthor.AuthorType),
		AvatarURL:        avatarURL,
		Rate:             rateValue,
		CreatedAt:        dbAuthor.CreatedAt,
		UpdatedAt:        dbAuthor.UpdatedAt,
	}
}
