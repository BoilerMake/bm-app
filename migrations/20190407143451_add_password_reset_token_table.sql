-- +goose Up
CREATE TABLE password_reset_tokens (
  id SERIAL UNIQUE PRIMARY KEY NOT NULL,
  uid INTEGER UNIQUE NOT NULL,
  token TEXT NOT NULL,
  valid_until TIMESTAMP default current_timestamp + interval '1 hour',
  FOREIGN KEY (uid) REFERENCES users (id)
);
-- +goose Down
DROP TABLE password_reset_tokens;