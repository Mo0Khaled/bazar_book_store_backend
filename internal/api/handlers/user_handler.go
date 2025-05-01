package handlers

import (
	"bazar_book_store/helpers"
	"bazar_book_store/internal/api/models"
	"bazar_book_store/internal/database"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
	"net/http"
)

func RegisterUserRoutes(r chi.Router) {
	r.Post("/register", Cfg.createUserHandler)
	r.Post("/login", Cfg.loginUserHandler)
	r.Post("/user/update-image", AuthMiddleware(Cfg.updateUserImageHandler))

}

func (apiCFG *ApiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	db := apiCFG.DB

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.RespondWithError(w, http.StatusForbidden, "Incorrect credentials")
		return
	}

	user, err := db.CreateUser(r.Context(), database.CreateUserParams{
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userToken, err := helpers.GenerateJWT(user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"user":    models.DBUserToUser(user),
		"token":   userToken,
		"message": "User created successfully",
	}
	helpers.RespondWithJSON(w, http.StatusCreated, response)
}

func (apiCFG *ApiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	db := apiCFG.DB

	user, err := db.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		helpers.RespondWithError(w, http.StatusForbidden, "Incorrect credentials")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password))
	if err != nil {
		helpers.RespondWithError(w, http.StatusForbidden, "Incorrect credentials")
		return
	}

	userToken, err := helpers.GenerateJWT(user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"user":    models.DBUserToUser(user),
		"token":   userToken,
		"message": "User signed in successfully",
	}
	helpers.RespondWithJSON(w, http.StatusOK, response)
}

func (apiCFG *ApiConfig) updateUserImageHandler(w http.ResponseWriter, r *http.Request, user database.User) {
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
	db := apiCFG.DB

	err = db.UpdateUserImage(r.Context(), database.UpdateUserImageParams{
		ID:        user.ID,
		AvatarUrl: helpers.ToNullString(&imageURL),
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Couldn't update user image")
		return
	}

	response := map[string]interface{}{
		"image_url": imageURL,
		"message":   "User image updated successfully",
	}
	helpers.RespondWithJSON(w, http.StatusOK, response)
}
