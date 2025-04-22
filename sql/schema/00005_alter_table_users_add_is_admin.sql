-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP
        COLUMN is_admin;
-- +goose StatementEnd
