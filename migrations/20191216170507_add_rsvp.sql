-- +goose Up
CREATE TABLE rsvps (
	id SERIAL UNIQUE NOT NULL,
	-- We'll mostly (only?) be getting RSVPs using a user's ID, so make
	-- that the primary key
	user_id INTEGER PRIMARY KEY REFERENCES users(id),

  will_attend BOOLEAN,
  accommodations TEXT,
  shirt_size TEXT,
  allergies TEXT,

  created_at TIMESTAMP default current_timestamp
);

-- +goose Down
DROP TABLE rsvps;
