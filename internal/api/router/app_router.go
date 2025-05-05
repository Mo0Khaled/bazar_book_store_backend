package router

import (
	"bazar_book_store/internal/api/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func InitRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"*"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           300,
			}),
	)
	v1Router := chi.NewRouter()
	handlers.RegisterUserRoutes(v1Router, h)
	handlers.RegisterAddressRoutes(v1Router)
	handlers.RegisterVendorsRoutes(v1Router)
	handlers.RegisterBooksRoutes(v1Router)
	handlers.RegisterCategoryRoutes(v1Router, h)
	handlers.RegisterAuthorRoutes(v1Router)
	handlers.RegisterImageHandlers(v1Router)
	handlers.RegisterApiTokensRoutes(v1Router)
	r.Mount("/v1", v1Router)

	return r
}
