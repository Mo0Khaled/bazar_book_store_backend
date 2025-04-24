-- +goose Up
-- +goose StatementBegin
CREATE TABLE book_favorites
(
    user_id    INT       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    book_id    INT       NOT NULL REFERENCES books (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, book_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE book_favorites;
-- +goose StatementEnd
