-- name: CreateAuthor :one

INSERT INTO authors (name, short_description, about, avatar_url, author_type, rate)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
