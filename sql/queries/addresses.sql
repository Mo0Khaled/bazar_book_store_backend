-- name: CreateAddress :one

INSERT INTO addresses
    (user_id, title, phone_number, governorate, city, address_details)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAddresses :many

SELECT *
FROM addresses
WHERE user_id = $1
ORDER BY id DESC;

-- name: UpdateAddress :one
UPDATE addresses
SET title           = COALESCE(@title, title),
    phone_number    = COALESCE(@phone_number, phone_number),
    governorate     = COALESCE(@governorate, governorate),
    city            = COALESCE(@city, city),
    address_details = COALESCE(@address_details, address_details),
    updated_at      = NOW()
WHERE id = @id
  AND user_id = @user_id
RETURNING *;

-- name: DeleteAddress :exec
DELETE
FROM addresses
WHERE id = $1
  AND user_id = $2;