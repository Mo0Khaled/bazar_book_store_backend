-- name: CreateCategory :one

INSERT INTO categories (name)
VALUES ($1)
RETURNING *;

-- name: GetCategoryByID :one

SELECT *
FROM categories
WHERE id = $1;