package models

import (
	"bazar_book_store/internal/database"
	"time"
)

type User struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarUrl *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DBUserToUser(dbUser database.User) User {
	var url *string
	if dbUser.AvatarUrl.Valid {
		url = &dbUser.AvatarUrl.String
	}
	return User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		AvatarUrl: url,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}
