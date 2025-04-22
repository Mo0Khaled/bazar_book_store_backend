package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {

	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errorResponse{Error: message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal json response %v with error %e\n", payload, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("failed to write response with error %e\n", err)
		return
	}
}

func DecodeBody[T any](w http.ResponseWriter, r *http.Request) (*T, bool) {
	decoder := json.NewDecoder(r.Body)
	var params T
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return nil, false
	}
	return &params, true
}
