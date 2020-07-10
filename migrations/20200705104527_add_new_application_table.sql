-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users
DROP COLUMN first_name,
DROP COLUMN last_name;

CREATE TABLE bm_applications (
	id SERIAL UNIQUE NOT NULL,
	-- We'll mostly (only?) be getting applications using a user's ID, so make
	-- that the primary key
	user_id INTEGER PRIMARY KEY REFERENCES users(id),
	decision SMALLINT,
	emailed_decision BOOLEAN,
	accepted_at TIMESTAMP,
	checked_in_at TIMESTAMP,

	rsvp BOOLEAN,

	school TEXT,
	gender TEXT,
	major TEXT,
	graduation_year TEXT,
	first_name TEXT,
	last_name TEXT,
	dietary_restrictions TEXT,
	github TEXT,
	resume_file TEXT,
	phone TEXT,
	is_first_hackathon BOOLEAN,
	referrer TEXT,
	race TEXT,
	why_bm TEXT,
	tac_18_or_older BOOLEAN,
	tac_mlh_code_of_conduct BOOLEAN,
	tac_mlh_contest_and_privacy BOOLEAN
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE bm_applications;