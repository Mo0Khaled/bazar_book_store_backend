package main

import (
	"bazar_book_store/internal/api/handlers"
	"bazar_book_store/internal/api/router"
	"bazar_book_store/internal/database"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Printf("Strting app")

	err := godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8081"
		log.Fatal("$PORT must be set")
	}

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("$DB_URL must be set")
	}
	connection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	handlers.Cfg = &handlers.ApiConfig{
		DB:  database.New(connection),
		RDB: rdb,
	}

	createdRouter := router.InitRouter(handlers.Cfg)
	srv := &http.Server{
		Handler: createdRouter,
		Addr:    ":" + portString,
	}
	log.Printf("Starting server on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
