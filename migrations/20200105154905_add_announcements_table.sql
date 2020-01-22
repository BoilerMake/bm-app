-- +goose Up
CREATE TABLE announcements (
  id SERIAL UNIQUE PRIMARY KEY NOT NULL,
  message TEXT NOT NULL,
  created_at TIMESTAMPTZ default current_timestamp
);
-- +goose Down
DROP TABLE announcements;
