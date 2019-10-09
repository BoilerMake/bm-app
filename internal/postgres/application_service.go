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
	Major                sql.NullString
	GraduationYear       sql.NullString
	ResumeFile           sql.NullString
	Gender               sql.NullString
	Race                 sql.NullString
	DietaryRestrictions  sql.NullString
	Github               sql.NullString
	Linkedin             sql.NullString
	WhyBM                sql.NullString
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
		Major:                a.Major.String,
		GraduationYear:       a.GraduationYear.String,
		ResumeFile:           a.ResumeFile.String,
		Gender:               a.Gender.String,
		Race:                 a.Race.String,
		DietaryRestrictions:  a.DietaryRestrictions.String,
		Github:               a.Github.String,
		Linkedin:             a.Linkedin.String,
		WhyBM:                a.WhyBM.String,
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
			major,
			graduation_year,
			resume_file,
			gender,
			race,
			dietary_restrictions,
			github,
			linkedin,
			why_bm,
			tac_18_or_older,
			tac_mlh_code_of_conduct,
			tac_mlh_contest_and_privacy
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`,
				newApp.UserID,
				newApp.School,
				newApp.Major,
				newApp.GraduationYear,
				newApp.ResumeFile,
				newApp.Gender,
				newApp.Race,
				newApp.DietaryRestrictions,
				newApp.Github,
				newApp.Linkedin,
				newApp.WhyBM,
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
				major = $2,
				graduation_year = $3,
				resume_file = $4,
				gender = $5,
				race = $6,
				dietary_restrictions = $7,
				github = $8,
				linkedin = $9,
				why_bm = $10
				tac_18_or_older = $11,
				tac_mlh_code_of_conduct = $12,
				tac_mlh_contest_and_privacy = $13
			WHERE user_id = $14`,
				newApp.School,
				newApp.Major,
				newApp.GraduationYear,
				newApp.ResumeFile,
				newApp.Gender,
				newApp.Race,
				newApp.DietaryRestrictions,
				newApp.Github,
				newApp.Linkedin,
				newApp.WhyBM,
				newApp.Is18OrOlder,
				newApp.MLHCodeOfConduct,
				newApp.MLHContestAndPrivacy,
				newApp.UserID,
			)
		} else {
			_, err = tx.Exec(`UPDATE bm7_applications
			SET
				school = $1,
				major = $2,
				graduation_year = $3,
				gender = $4,
				race = $5,
				dietary_restrictions = $6,
				github = $7,
				linkedin = $8,
				why_bm = $9,
				tac_18_or_older = $10,
				tac_mlh_code_of_conduct = $11,
				tac_mlh_contest_and_privacy = $12
			WHERE user_id = $13`,
				newApp.School,
				newApp.Major,
				newApp.GraduationYear,
				newApp.Gender,
				newApp.Race,
				newApp.DietaryRestrictions,
				newApp.Github,
				newApp.Linkedin,
				newApp.WhyBM,
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
			major,
			graduation_year,
			resume_file,
			gender,
			race,
			dietary_restrictions,
			github,
			linkedin,
			why_bm,
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
		&dba.Major,
		&dba.GraduationYear,
		&dba.ResumeFile,
		&dba.Gender,
		&dba.Race,
		&dba.DietaryRestrictions,
		&dba.Github,
		&dba.Linkedin,
		&dba.WhyBM,
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
