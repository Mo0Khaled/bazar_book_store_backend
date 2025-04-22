package models

import (
	"bazar_book_store/internal/database"
	"time"
)

type Address struct {
	ID             int32     `json:"id"`
	Title          string    `json:"title"`
	AddressDetails string    `json:"address_details"`
	Governorate    string    `json:"governorate"`
	PhoneNumber    string    `json:"phone_number"`
	City           string    `json:"city"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func DBAddressToAddress(dbUser database.Address) Address {
	return Address{
		ID:             dbUser.ID,
		Title:          dbUser.Title,
		AddressDetails: dbUser.AddressDetails,
		Governorate:    dbUser.Governorate,
		PhoneNumber:    dbUser.PhoneNumber,
		City:           dbUser.City,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
	}
}

func DBAddressesToAddresses(dbAddresses []database.Address) []Address {
	addresses := make([]Address, len(dbAddresses))
	for i, dbAddress := range dbAddresses {
		addresses[i] = DBAddressToAddress(dbAddress)
	}
	return addresses
}
