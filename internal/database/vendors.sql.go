// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: vendors.sql

package database

import (
	"context"
)

const createVendor = `-- name: CreateVendor :one

INSERT INTO vendors (name, avatar_url, rate)
VALUES ($1, $2, $3)
RETURNING id, name, avatar_url, rate, created_at, updated_at
`

type CreateVendorParams struct {
	Name      string
	AvatarUrl string
	Rate      string
}

func (q *Queries) CreateVendor(ctx context.Context, arg CreateVendorParams) (Vendor, error) {
	row := q.db.QueryRowContext(ctx, createVendor, arg.Name, arg.AvatarUrl, arg.Rate)
	var i Vendor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.AvatarUrl,
		&i.Rate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getVendors = `-- name: GetVendors :many

SELECT id, name, avatar_url, rate, created_at, updated_at FROM vendors
`

func (q *Queries) GetVendors(ctx context.Context) ([]Vendor, error) {
	rows, err := q.db.QueryContext(ctx, getVendors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Vendor
	for rows.Next() {
		var i Vendor
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.AvatarUrl,
			&i.Rate,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
