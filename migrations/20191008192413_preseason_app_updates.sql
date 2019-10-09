-- +goose Up
ALTER TABLE bm7_applications
DROP COLUMN is_first_hackathon,
DROP COLUMN project_idea,
DROP COLUMN team_members,
ADD COLUMN why_bm TEXT;

-- +goose Down
ALTER TABLE bm7_applications
ADD COLUMN is_first_hackathon BOOLEAN,
ADD COLUMN project_idea TEXT,
ADD COLUMN team_members TEXT[3];
