-- +goose Up
ALTER TABLE users
DROP COLUMN project_idea,
DROP COLUMN team_members;

CREATE TABLE bm7_applications (
	id SERIAL UNIQUE NOT NULL,
	-- We'll mostly (only?) be getting applications using a user's ID, so make
	-- that the primary key
	user_id INTEGER PRIMARY KEY REFERENCES users(id),
	decision SMALLINT,
	emailed_decision BOOLEAN,
	checked_in_at TIMESTAMP,

	rsvp BOOLEAN,

	school TEXT,
	gender TEXT,
	major TEXT,
	graduation_year TEXT,
	dietary_restrictions TEXT,
	github TEXT,
	linkedin TEXT,
	resume_file TEXT,
	is_first_hackathon BOOLEAN,
	race TEXT,
	shirt_size TEXT,
	project_idea TEXT,
	team_members TEXT[3],
	tac_18_or_older BOOLEAN,
	tac_mlh_code_of_conduct BOOLEAN,
	tac_mlh_contest_and_privacy BOOLEAN
);

-- +goose Down
ALTER TABLE users
ADD COLUMN project_idea TEXT,
ADD COLUMN team_members TEXT[3];

DROP TABLE bm7_applications;
