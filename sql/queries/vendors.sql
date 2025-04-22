-- name: CreateVendor :one

INSERT INTO vendors (name, avatar_url, rate)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVendors :many

SELECT * FROM vendors;