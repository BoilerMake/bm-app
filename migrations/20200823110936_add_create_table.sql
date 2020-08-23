-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE raffles (
    code TEXT PRIMARY KEY,
    start_time BIGINT,
    end_time BIGINT,
    points  INTEGER
);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE raffles;