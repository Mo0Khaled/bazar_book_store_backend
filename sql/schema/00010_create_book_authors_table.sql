-- +goose Up
-- +goose StatementBegin
CREATE TABLE book_authors
(
    book_id   INT NOT NULL REFERENCES books (id) ON DELETE CASCADE,
    author_id INT NOT NULL REFERENCES authors (id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE book_authors;
-- +goose StatementEnd
