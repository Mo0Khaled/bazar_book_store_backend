// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: authors.sql

package database

import (
	"context"
	"database/sql"
)

const createAuthor = `-- name: CreateAuthor :one

INSERT INTO authors (name, short_description, about, avatar_url, author_type, rate)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, short_description, about, avatar_url, author_type, rate, created_at, updated_at
`

type CreateAuthorParams struct {
	Name             string
	ShortDescription string
	About            string
	AvatarUrl        sql.NullString
	AuthorType       AuthorTypeEnum
	Rate             string
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	row := q.db.QueryRowContext(ctx, createAuthor,
		arg.Name,
		arg.ShortDescription,
		arg.About,
		arg.AvatarUrl,
		arg.AuthorType,
		arg.Rate,
	)
	var i Author
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ShortDescription,
		&i.About,
		&i.AvatarUrl,
		&i.AuthorType,
		&i.Rate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
