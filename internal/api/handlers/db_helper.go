package handlers

import (
	"bazar_book_store/internal/database"
	"github.com/redis/go-redis/v9"
)

type ApiConfig struct {
	DB  *database.Queries
	RDB *redis.Client
}

var Cfg *ApiConfig
