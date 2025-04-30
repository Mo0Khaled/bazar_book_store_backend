package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/database"
	"github.com/go-chi/chi/v5"
	"mime/multipart"
	"net/http"
)

func RegisterImageHandlers(r chi.Router) {
	r.Post("/upload_image", AuthMiddleware(Cfg.uploadImageHandler))
}

func (apiCFG *ApiConfig) uploadImageHandler(w http.ResponseWriter, r *http.Request, _ database.User) {
	file, _, err := r.FormFile("image")
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Image isn't exists")
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	imageURL, err := helpers.UploadImage(file)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not upload the image")
		return
	}

	response := map[string]interface{}{
		"image_url": imageURL,
		"message":   "Image uploaded successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}
