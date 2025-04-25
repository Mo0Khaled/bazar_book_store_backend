-- +goose Up
-- +goose StatementBegin
ALTER TABLE books
    ADD COLUMN avatar_url TEXT NOT NULL default '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE books DROP COLUMN avatar_url;
-- +goose StatementEnd
