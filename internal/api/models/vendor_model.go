package models

import (
	"bazar_book_store/internal/database"
	"strconv"
	"time"
)

type Vendor struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	Rate      float64   `json:"rate"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DBVendorToVendor(dbVendor database.Vendor) Vendor {
	floatValue, err := strconv.ParseFloat(dbVendor.Rate, 64)
	if err != nil {
		floatValue = 0.1
	}
	return Vendor{
		ID:        dbVendor.ID,
		Name:      dbVendor.Name,
		AvatarURL: dbVendor.AvatarUrl,
		Rate:      floatValue,
		CreatedAt: dbVendor.CreatedAt,
		UpdatedAt: dbVendor.UpdatedAt,
	}
}

func DBVendorsToVendors(dbVendors []database.Vendor) []Vendor {
	vendors := make([]Vendor, len(dbVendors))

	for i, dbVendor := range dbVendors {
		vendors[i] = DBVendorToVendor(dbVendor)
	}
	return vendors
}
