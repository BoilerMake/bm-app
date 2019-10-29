-- +goose Up
ALTER TABLE users 
DROP COLUMN phone;

ALTER TABLE bm7_applications
DROP COLUMN linkedin,
ADD COLUMN is_first_hackathon BOOLEAN,
ADD COLUMN referrer TEXT,
ADD COLUMN phone TEXT;

-- +goose Down
ALTER TABLE users
ADD COLUMN phone TEXT;

ALTER TABLE bm7_applications
ADD COLUMN linkedin TEXT,
DROP COLUMN is_first_hackathon,
DROP COLUMN referrer,
DROP COLUMN phone;
