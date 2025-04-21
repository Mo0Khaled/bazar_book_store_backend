package main

import (
	"bazar_book_store/internal/api/handlers"
	"bazar_book_store/internal/api/router"
	"bazar_book_store/internal/database"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	portString := os.Getenv("PORT")
	if portString == "" {
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
	handlers.Cfg = &handlers.ApiConfig{
		DB: database.New(connection),
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
