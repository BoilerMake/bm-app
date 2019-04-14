-- +goose Up
CREATE TABLE password_reset_tokens (
  id SERIAL UNIQUE PRIMARY KEY NOT NULL,
  uid INTEGER NOT NULL,
  tokenID TEXT UNIQUE NOT NULL,
  hashedToken TEXT UNIQUE NOT NULL,
  created_at TIMESTAMPTZ default current_timestamp,
  FOREIGN KEY (uid) REFERENCES users (id)
);
-- +goose Down
DROP TABLE password_reset_tokens;