package handlers

import (
	"bazar_book_store/internal/database"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDB struct{}

func (m *mockDB) AddBookAuthor(ctx context.Context, arg database.AddBookAuthorParams) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) AddBookCategory(ctx context.Context, arg database.AddBookCategoryParams) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) AddBookFavorite(ctx context.Context, arg database.AddBookFavoriteParams) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) CountBooks(ctx context.Context) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) CreateAddress(ctx context.Context, arg database.CreateAddressParams) (database.Address, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) CreateApiToken(ctx context.Context, arg database.CreateApiTokenParams) (database.ApiToken, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) CreateAuthor(ctx context.Context, arg database.CreateAuthorParams) (database.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) CreateBook(ctx context.Context, arg database.CreateBookParams) (database.Book, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) CreateVendor(ctx context.Context, arg database.CreateVendorParams) (database.Vendor, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) DeleteAddress(ctx context.Context, arg database.DeleteAddressParams) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) GetAddresses(ctx context.Context, userID int32) ([]database.Address, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) GetApiToken(ctx context.Context, apiToken string) (database.ApiToken, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) GetBooks(ctx context.Context, arg database.GetBooksParams) ([]database.Book, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) GetBooksDetails(ctx context.Context, arg database.GetBooksDetailsParams) ([]database.GetBooksDetailsRow, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) GetCategoryByID(ctx context.Context, id int32) (database.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) GetFavoriteBooks(ctx context.Context, userID int32) ([]database.GetFavoriteBooksRow, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) GetUser(ctx context.Context, id int32) (database.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) GetVendors(ctx context.Context) ([]database.Vendor, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) RemoveBookFavorite(ctx context.Context, arg database.RemoveBookFavoriteParams) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) UpdateAddress(ctx context.Context, arg database.UpdateAddressParams) (database.Address, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) UpdateUserImage(ctx context.Context, arg database.UpdateUserImageParams) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockDB) CreateCategory(ctx context.Context, name string) (database.Category, error) {
	if name == "fail" {
		return database.Category{}, errors.New("failed to create category")
	}
	return database.Category{ID: 1, Name: name}, nil
}

func TestCreateCategoryHandler_Success(t *testing.T) {
	mock := &mockDB{}
	cfg := &ApiConfig{DB: mock}
	h := NewHandler(cfg)

	body := map[string]string{"name": "Fiction"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(jsonBody))

	rr := httptest.NewRecorder()

	dummyUser := database.User{ID: 1, Email: "test@example.com"}

	h.createCategoryHandler(rr, req, dummyUser)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "Fiction")
}
