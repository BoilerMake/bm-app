-- +goose Up
CREATE TABLE users (
  id SERIAL UNIQUE PRIMARY KEY NOT NULL,
  role SMALLINT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  phone TEXT,

  project_idea TEXT,
  team_members TEXT[3],

  created_at TIMESTAMP default current_timestamp
);

-- +goose Down
DROP TABLE users
