-- +goose Up
ALTER TABLE users
ADD COLUMN is_active boolean,
ADD COLUMN confirmation_code text;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_active,
DROP COLUMN confirmation_code;
