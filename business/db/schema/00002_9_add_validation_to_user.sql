-- Description: Add new fields validation_token and validated.

-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN validation_token TEXT NULL,
    ADD COLUMN validated BOOL DEFAULT FALSE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN validation_token,
    DROP COLUMN validated;
-- +goose StatementEnd