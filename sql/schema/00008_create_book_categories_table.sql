-- +goose Up
-- +goose StatementBegin
CREATE TABLE book_categories
(
    book_id     INT NOT NULL REFERENCES books (id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE book_categories;
-- +goose StatementEnd
