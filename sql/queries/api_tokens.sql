-- name: CreateApiToken :one

INSERT INTO api_tokens (api_token, expires_at, request_limit, last_reset)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: GetApiToken :one

SELECT *
FROM api_tokens
WHERE api_token = $1;