-- +goose Up
ALTER TABLE bm_applications
ADD COLUMN check_in_status BOOLEAN;

-- +goose Down
ALTER TABLE bm_applications
DROP COLUMN check_in_status;