-- +goose Up
ALTER TABLE bm7_applications
ADD COLUMN accepted_at TIMESTAMP;

-- +goose Down
ALTER TABLE bm7_applications
DROP COLUMN accepted_at;
