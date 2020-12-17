-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE bm_applications
ADD COLUMN check_in_status BOOLEAN;
-- +goose StatementEnd
​
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE bm_applications
DROP COLUMN check_in_status;
-- +goose StatementEnd