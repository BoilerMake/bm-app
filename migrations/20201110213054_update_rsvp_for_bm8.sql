-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE rsvps
DROP COLUMN accommodations,
DROP COLUMN allergies,
ADD COLUMN oncampus BOOLEAN,
ADD COLUMN shipping_addr TEXT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE rsvps TEXT
ADD COLUMN accommodations TEXT,
ADD COLUMN allergies,
DROP COLUMN oncampus,
DROP COLUMN shipping_addr;
