-- +goose Up
-- +goose StatementBegin
CREATE TABLE api_tokens
(
    api_token     TEXT UNIQUE NOT NULL PRIMARY KEY,
    created_at    TIMESTAMP DEFAULT now(),
    expires_at    TIMESTAMP   NOT NULL,
    request_limit INT         NOT NULL,
    last_reset    TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE api_tokens;
-- +goose StatementEnd
