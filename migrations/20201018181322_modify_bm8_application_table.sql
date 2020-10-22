-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE bm_applications
ADD COLUMN location TEXT,
ADD COLUMN other_school TEXT,
ADD COLUMN other_major TEXT,
ADD COLUMN proj_idea TEXT,
DROP COLUMN referrer,
DROP COLUMN dietary_restrictions,
DROP COLUMN race;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE bm_applications
DROP COLUMN location,
DROP COLUMN other_school TEXT,
DROP COLUMN other_major,
DROP COLUMN proj_idea,
ADD COLUMN referrer TEXT,
ADD COLUMN dietary_restrictions TEXT,
ADD COLUMN race TEXT;