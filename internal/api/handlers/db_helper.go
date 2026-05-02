package handlers

import (
	"bazar_book_store/internal/database"
	"github.com/redis/go-redis/v9"
)

type ApiConfig struct {
	DB  database.Querier
	RDB *redis.Client
}

var Cfg *ApiConfig

type Handler struct {
	Cfg *ApiConfig
}

func NewHandler(cfg *ApiConfig) *Handler {
	return &Handler{Cfg: cfg}
}
