-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE bm7_applications
ADD COLUMN first_name TEXT NOT NULL,
ADD COLUMN last_name TEXT NOT NULL;

ALTER TABLE users
DROP COLUMN first_name,
DROP COLUMN last_name;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE bm7_applications
DROP COLUMN first_name,
DROP COLUMN last_name;

ALTER TABLE users
ADD COLUMN first_name TEXT NOT NULL,
ADD COLUMN last_name TEXT NOT NULL;