package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/api/models"
	"bazar_book_store/internal/database"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func RegisterAddressRoutes(r chi.Router) {
	r.Post("/addresses", AuthMiddleware(Cfg.createAddressHandler))
	r.Get("/addresses", AuthMiddleware(Cfg.getAddressesHandler))
	r.Put("/addresses", AuthMiddleware(Cfg.updateAddressHandler))
	r.Delete("/addresses/{addressID}", AuthMiddleware(Cfg.deleteAddressHandler))

}

func (apiCFG *ApiConfig) createAddressHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title          string `json:"title"`
		PhoneNumber    string `json:"phone_number"`
		Governorate    string `json:"governorate"`
		City           string `json:"city"`
		AddressDetails string `json:"address_details"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	if params.Title == "" || params.City == "" || params.AddressDetails == "" || params.PhoneNumber == "" || params.Governorate == "" {
		helpers.RespondWithError(w, http.StatusForbidden, "all fields are required")
		return
	}

	db := apiCFG.DB

	address, err := db.CreateAddress(r.Context(), database.CreateAddressParams{
		UserID:         user.ID,
		Title:          params.Title,
		PhoneNumber:    params.PhoneNumber,
		Governorate:    params.Governorate,
		City:           params.City,
		AddressDetails: params.AddressDetails,
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not create address")
		return
	}

	response := map[string]interface{}{
		"address": models.DBAddressToAddress(address),
		"message": "Address created successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}

func (apiCFG *ApiConfig) getAddressesHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	db := apiCFG.DB

	addresses, err := db.GetAddresses(r.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not get addresses")
		return
	}

	response := map[string]interface{}{
		"addresses": models.DBAddressesToAddresses(addresses),
		"message":   "Successfully fetched addresses",
	}
	helpers.RespondWithJSON(w, http.StatusOK, response)
}

func (apiCFG *ApiConfig) updateAddressHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		ID             int32  `json:"id"`
		Title          string `json:"title"`
		PhoneNumber    string `json:"phone_number"`
		Governorate    string `json:"governorate"`
		City           string `json:"city"`
		AddressDetails string `json:"address_details"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	db := apiCFG.DB

	address, err := db.UpdateAddress(r.Context(), database.UpdateAddressParams{
		ID:             params.ID,
		UserID:         user.ID,
		Title:          params.Title,
		PhoneNumber:    params.PhoneNumber,
		Governorate:    params.Governorate,
		City:           params.City,
		AddressDetails: params.AddressDetails,
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not update address")
		return
	}

	response := map[string]interface{}{
		"address": models.DBAddressToAddress(address),
		"message": "Address updated successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}

func (apiCFG *ApiConfig) deleteAddressHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	addressIDStr := chi.URLParam(r, "addressID")

	addressID, err := strconv.Atoi(addressIDStr)

	if err != nil {
		helpers.RespondWithError(w, http.StatusForbidden, "wrong address ID")
	}

	db := apiCFG.DB

	err = db.DeleteAddress(r.Context(), database.DeleteAddressParams{
		ID:     int32(addressID),
		UserID: user.ID,
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not delete address")
		return
	}

	response := map[string]interface{}{
		"message": "Address deleted successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}
