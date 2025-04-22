-- name: CreateUser :one

INSERT INTO users
    (id, name, email, password_hash, avatar_url, is_admin)
VALUES (DEFAULT, $1, $2, $3, $4,$5)
RETURNING *;

-- name: GetUser :one

SELECT *
FROM users
where id = $1;

-- name: GetUserByEmail :one

SELECT *
FROM users
where email = $1;
