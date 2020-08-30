-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION hstore;

ALTER TABLE bm_applications
ADD COLUMN points INTEGER,
ADD COLUMN raffles_claimed hstore;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE bm_applications
DROP COLUMN points,
DROP COLUMN raffles_claimed;