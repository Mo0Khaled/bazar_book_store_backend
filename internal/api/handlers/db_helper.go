package handlers

import "bazar_book_store/internal/database"

type ApiConfig struct {
	DB *database.Queries
}

var Cfg *ApiConfig
