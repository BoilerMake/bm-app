-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE bm_applications
ADD COLUMN points INTEGER;

CREATE TABLE raffles (
    raffle_id TEXT PRIMARY KEY,
    start_time BIGINT,
    end_time BIGINT,
    points INTEGER
);

CREATE TABLE raffle_hacker (
    user_id INTEGER REFERENCES bm_applications(user_id),
    raffle_id TEXT REFERENCES raffles(raffle_id),
    PRIMARY KEY(user_id, raffle_id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE bm_applications
DROP COLUMN points;

DROP TABLE raffles;

DROP TABLE raffle_hacker;
