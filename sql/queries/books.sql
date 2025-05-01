-- name: CreateBook :one

INSERT INTO books (vendor_id, title, description, price, rate, avatar_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: AddBookCategory :exec
INSERT INTO book_categories (book_id, category_id)
VALUES ($1, $2);

-- name: AddBookAuthor :exec
INSERT INTO book_authors (book_id, author_id)
VALUES ($1, $2);

-- name: AddBookFavorite :exec
INSERT INTO book_favorites (user_id, book_id)
VALUES ($1, $2);

-- name: RemoveBookFavorite :exec
DELETE
FROM book_favorites
WHERE user_id = $1
  AND book_id = $2;

-- name: GetBooks :many
SELECT *
FROM books
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: CountBooks :one
SELECT COUNT(*)
FROM books;


-- name: GetBooksDetails :many

SELECT b.id                                                      AS book_id,
       b.title,
       b.description,
       b.avatar_url                                              AS book_avatar_url,
       b.price,
       b.rate,
       b.vendor_id,
       CASE WHEN bf.user_id IS NOT NULL THEN true ELSE false END AS is_favorite,
       b.created_at,
       b.updated_at,


       v.id                                                      AS vendor_id,
       v.name                                                    AS vendor_name,
       v.avatar_url                                              AS vendor_avatar_url,
       v.rate                                                    AS vendor_rate,

       a.id                                                      AS author_id,
       a.name                                                    AS author_name,
       a.short_description                                       AS author_short_description,
       a.about                                                   AS author_about,
       a.avatar_url                                              AS author_avatar_url,
       a.rate                                                    AS author_rate,
       a.author_type,

       c.id                                                      AS category_id,
       c.name                                                    AS category_name

FROM books b
         JOIN vendors v ON b.vendor_id = v.id
         LEFT JOIN book_authors ba ON b.id = ba.book_id
         LEFT JOIN authors a ON a.id = ba.author_id
         LEFT JOIN book_categories bc ON b.id = bc.book_id
         LEFT JOIN categories c ON c.id = bc.category_id
         LEFT JOIN book_favorites bf
                   ON b.id = bf.book_id
                       AND user_id = $1
WHERE (sqlc.narg(category_id)::int IS NULL OR c.id = sqlc.narg(category_id))
  AND (sqlc.narg(vendor_id)::int IS NULL OR b.vendor_id = sqlc.narg(vendor_id))
  AND (sqlc.narg(author_id)::int IS NULL OR a.id = sqlc.narg(author_id))
  AND (sqlc.narg(book_id)::int IS NULL OR b.id = sqlc.narg(book_id))
LIMIT $2 OFFSET $3;

-- name: GetFavoriteBooks :many

SELECT id, title, price, avatar_url
From books
         JOIN book_favorites ON book_favorites.user_id = $1 AND books.id = book_favorites.book_id
ORDER BY book_favorites.created_at DESC;


