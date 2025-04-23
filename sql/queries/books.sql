-- name: CreateBook :one

INSERT INTO books (vendor_id, title, description, price, rate)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: AddBookCategory :exec
INSERT INTO book_categories (book_id, category_id)
VALUES ($1, $2);

-- name: AddBookAuthor :exec
INSERT INTO book_authors (book_id, author_id)
VALUES ($1, $2);

-- name: GetBooks :many

SELECT *
FROM books
ORDER BY id DESC;