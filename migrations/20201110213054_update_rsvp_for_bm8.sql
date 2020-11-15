-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE rsvps
DROP COLUMN accommodations,
DROP COLUMN allergies,
ADD COLUMN oncampus BOOLEAN,
ADD COLUMN street_address TEXT,
ADD COLUMN city TEXT,
ADD COLUMN state TEXT,
ADD COLUMN country TEXT,
ADD COLUMN zip_code TEXT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE rsvps TEXT
ADD COLUMN accommodations TEXT,
ADD COLUMN allergies,
DROP COLUMN oncampus,
DROP COLUMN street_address,
DROP COLUMN city,
DROP COLUMN country,
DROP COLUMN zip_code;
