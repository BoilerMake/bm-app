package postgres

import (
	"database/sql"

	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/lib/pq"
)

// UserService is a PostgreSQL implementation of models.UserService
type ApplicationService struct {
	DB *sql.DB
}

// A dbApplication is like a models.Application, but it can read in null fields
// from the database without panicking.
type dbApplication struct {
	ID                   sql.NullInt64
	UserID               sql.NullInt64
	Decision             sql.NullInt64
	EmailedDecision      sql.NullBool
	CheckedInAt          pq.NullTime
	RSVP                 sql.NullBool
	School               sql.NullString
	Gender               sql.NullString
	Major                sql.NullString
	GraduationYear       sql.NullString
	DietaryRestrictions  sql.NullString
	Github               sql.NullString
	Linkedin             sql.NullString
	ResumeFile           sql.NullString
	IsFirstHackathon     sql.NullBool
	Race                 sql.NullString
	ShirtSize            sql.NullString
	ProjectIdea          sql.NullString
	TeamMembers          []string
	Is18OrOlder          sql.NullBool
	MLHCodeOfConduct     sql.NullBool
	MLHContestAndPrivacy sql.NullBool
}

// toModel converts a database specific dbApplication to the more generic
// Application struct.
func (a *dbApplication) toModel() *models.Application {
	return &models.Application{
		ID:                   int(a.ID.Int64),
		UserID:               int(a.UserID.Int64),
		Decision:             int(a.Decision.Int64),
		EmailedDecision:      a.EmailedDecision.Bool,
		CheckedInAt:          a.CheckedInAt.Time,
		RSVP:                 a.RSVP.Bool,
		School:               a.School.String,
		Gender:               a.Gender.String,
		Major:                a.Major.String,
		GraduationYear:       a.GraduationYear.String,
		DietaryRestrictions:  a.DietaryRestrictions.String,
		Github:               a.Github.String,
		Linkedin:             a.Linkedin.String,
		ResumeFile:           a.ResumeFile.String,
		IsFirstHackathon:     a.IsFirstHackathon.Bool,
		Race:                 a.Race.String,
		ShirtSize:            a.ShirtSize.String,
		ProjectIdea:          a.ProjectIdea.String,
		TeamMembers:          a.TeamMembers,
		Is18OrOlder:          a.Is18OrOlder.Bool,
		MLHCodeOfConduct:     a.MLHCodeOfConduct.Bool,
		MLHContestAndPrivacy: a.MLHContestAndPrivacy.Bool,
	}
}

// CreateOrUpdate tries to make a new Application, if one already exists then
// it updates the existing one.
func (s *ApplicationService) CreateOrUpdate(newApp *models.Application) (err error) {
	err = newApp.Validate()
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	oldApp, err := s.GetByUserID(newApp.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			// Application hasn't been made yet
			// Validate resume was uploaded
			if newApp.ResumeFile == "" {
				return models.ErrMissingResume
			}
			_, err = tx.Exec(`INSERT INTO bm7_applications (
			user_id,
			school,
			gender,
			major,
			graduation_year,
			dietary_restrictions,
			github,
			linkedin,
			resume_file,
			is_first_hackathon,
			race,
			shirt_size,
			project_idea,
			team_members,
			tac_18_or_older,
			tac_mlh_code_of_conduct,
			tac_mlh_contest_and_privacy
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17);`,
				newApp.UserID,
				newApp.School,
				newApp.Gender,
				newApp.Major,
				newApp.GraduationYear,
				newApp.DietaryRestrictions,
				newApp.Github,
				newApp.Linkedin,
				newApp.ResumeFile,
				newApp.IsFirstHackathon,
				newApp.Race,
				newApp.ShirtSize,
				newApp.ProjectIdea,
				pq.Array(newApp.TeamMembers),
				newApp.Is18OrOlder,
				newApp.MLHCodeOfConduct,
				newApp.MLHContestAndPrivacy,
			)

			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return rollbackErr
				}
				return err
			}
		} else {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	} else {
		// Application already exists, so update it

		// Validate resume was uploaded or already exists
		if newApp.ResumeFile == "" && oldApp.ResumeFile == "" {
			return models.ErrMissingResume
		}

		if newApp.ResumeFile != "" {
			_, err = tx.Exec(`UPDATE bm7_applications
			SET
				school = $1,
				gender = $2,
				major = $3,
				graduation_year = $4,
				dietary_restrictions = $5,
				github = $6,
				linkedin = $7,
				resume_file = $8,
				is_first_hackathon = $9,
				race = $10,
				shirt_size = $11,
				project_idea = $12,
				team_members = $13,
				tac_18_or_older = $14,
				tac_mlh_code_of_conduct = $15,
				tac_mlh_contest_and_privacy = $16
			WHERE user_id = $17`,
				newApp.School,
				newApp.Gender,
				newApp.Major,
				newApp.GraduationYear,
				newApp.DietaryRestrictions,
				newApp.Github,
				newApp.Linkedin,
				newApp.ResumeFile,
				newApp.IsFirstHackathon,
				newApp.Race,
				newApp.ShirtSize,
				newApp.ProjectIdea,
				pq.Array(newApp.TeamMembers),
				newApp.Is18OrOlder,
				newApp.MLHCodeOfConduct,
				newApp.MLHContestAndPrivacy,
				newApp.UserID,
			)
		} else {
			_, err = tx.Exec(`UPDATE bm7_applications
			SET
				school = $1,
				gender = $2,
				major = $3,
				graduation_year = $4,
				dietary_restrictions = $5,
				github = $6,
				linkedin = $7,
				is_first_hackathon = $8,
				race = $9,
				shirt_size = $10,
				project_idea = $11,
				team_members = $12,
				tac_18_or_older = $13,
				tac_mlh_code_of_conduct = $14,
				tac_mlh_contest_and_privacy = $15
			WHERE user_id = $16`,
				newApp.School,
				newApp.Gender,
				newApp.Major,
				newApp.GraduationYear,
				newApp.DietaryRestrictions,
				newApp.Github,
				newApp.Linkedin,
				newApp.IsFirstHackathon,
				newApp.Race,
				newApp.ShirtSize,
				newApp.ProjectIdea,
				pq.Array(newApp.TeamMembers),
				newApp.Is18OrOlder,
				newApp.MLHCodeOfConduct,
				newApp.MLHContestAndPrivacy,
				newApp.UserID,
			)
		}

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	err = tx.Commit()
	return err
}

// GetByUserID returns a single Application with the given user id.
func (s *ApplicationService) GetByUserID(uid int) (*models.Application, error) {
	var dba dbApplication

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(`SELECT
			id,
			user_id,
			decision,
			rsvp,
			school,
			gender,
			major,
			graduation_year,
			dietary_restrictions,
			github,
			linkedin,
			resume_file,
			is_first_hackathon,
			race,
			shirt_size,
			project_idea,
			team_members,
			tac_18_or_older,
			tac_mlh_code_of_conduct,
			tac_mlh_contest_and_privacy
		FROM bm7_applications
		WHERE user_id = $1`, uid).Scan(
		&dba.ID,
		&dba.UserID,
		&dba.Decision,
		&dba.RSVP,
		&dba.School,
		&dba.Gender,
		&dba.Major,
		&dba.GraduationYear,
		&dba.DietaryRestrictions,
		&dba.Github,
		&dba.Linkedin,
		&dba.ResumeFile,
		&dba.IsFirstHackathon,
		&dba.Race,
		&dba.ShirtSize,
		&dba.ProjectIdea,
		pq.Array(&dba.TeamMembers),
		&dba.Is18OrOlder,
		&dba.MLHCodeOfConduct,
		&dba.MLHContestAndPrivacy,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	err = tx.Commit()
	return dba.toModel(), err
}
