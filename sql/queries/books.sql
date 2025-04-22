-- name: CreateBook :one

INSERT INTO books (vendor_id, title, description, price, rate)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: AddBookCategory :exec
INSERT INTO book_categories (book_id, category_id)
VALUES ($1, $2);