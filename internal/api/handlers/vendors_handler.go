package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/api/models"
	"bazar_book_store/internal/database"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"net/http"
	"strconv"
)

func RegisterVendorsRoutes(r chi.Router) {
	r.Post("/vendors", AdminOnlyMiddleware(Cfg.createVendorHandler))
	r.Get("/vendors", Cfg.getVendorsHandler)

}

func (apiCFG *ApiConfig) createVendorHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name      string  `json:"name"`
		AvatarURL string  `json:"avatar_url"`
		Rate      float64 `json:"rate"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	db := apiCFG.DB

	vendor, err := db.CreateVendor(r.Context(), database.CreateVendorParams{
		Name:      params.Name,
		AvatarUrl: params.AvatarURL,
		Rate:      strconv.FormatFloat(params.Rate, 'f', 2, 64),
	})

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				if pqErr.Constraint == "vendors_name_key" {
					helpers.RespondWithError(w, http.StatusBadRequest, "Vendor name already exists")
					return
				}
			case "check_violation":
				if pqErr.Constraint == "vendors_rate_check" {
					helpers.RespondWithError(w, http.StatusBadRequest, "Rate must be between 1.0 and 5.0")
					return
				}
			}
		}
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not create vendor")
		return
	}

	response := map[string]interface{}{
		"vendor":  models.DBVendorToVendor(vendor),
		"message": "Vendor created successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}

func (apiCFG *ApiConfig) getVendorsHandler(w http.ResponseWriter, r *http.Request) {
	db := apiCFG.DB

	vendors, err := db.GetVendors(r.Context())

	if err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, "Couldn't fetch vendors")
		return
	}

	response := map[string]interface{}{
		"vendors": models.DBVendorsToVendors(vendors),
		"message": "Vendors fetched successfully",
	}
	helpers.RespondWithJSON(w, http.StatusOK, response)
}
