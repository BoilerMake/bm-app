package postgres

import (
	"database/sql"

	"github.com/BoilerMake/bm-app/internal/models"
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
	AcceptedAt           pq.NullTime
	CheckedInAt          pq.NullTime
	RSVP                 sql.NullBool
	School               sql.NullString
	Major                sql.NullString
	GraduationYear       sql.NullString
	FirstName            sql.NullString
	LastName             sql.NullString
	ResumeFile           sql.NullString
	Phone                sql.NullString
	Gender               sql.NullString
	Github               sql.NullString
	Location			 sql.NullString
	OtherMajor			 sql.NullString
	IsFirstHackathon     sql.NullBool
	WhyBM                sql.NullString
	ProjIdea			 sql.NullString
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
		AcceptedAt:           a.AcceptedAt.Time,
		CheckedInAt:          a.CheckedInAt.Time,
		RSVP:                 a.RSVP.Bool,
		School:               a.School.String,
		Major:                a.Major.String,
		GraduationYear:       a.GraduationYear.String,
		FirstName:            a.FirstName.String,
		LastName:             a.LastName.String,
		ResumeFile:           a.ResumeFile.String,
		Phone:                a.Phone.String,
		Gender:               a.Gender.String,
		Github:               a.Github.String,
		Location:			  a.Location.String,
		OtherMajor:			  a.OtherMajor.String,
		IsFirstHackathon:     a.IsFirstHackathon.Bool,
		WhyBM:                a.WhyBM.String,
		ProjIdea:			  a.ProjIdea.String,
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
			_, err = tx.Exec(`INSERT INTO bm_applications (
			user_id,
			school,
			major,
			graduation_year,
			first_name,
			last_name,
			resume_file,
			phone,
			gender,
			github,
			location,
			other_major,
			is_first_hackathon,
			why_bm,
			proj_idea,
			tac_18_or_older,
			tac_mlh_code_of_conduct,
			tac_mlh_contest_and_privacy
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);`,
				newApp.UserID,
				newApp.School,
				newApp.Major,
				newApp.GraduationYear,
				newApp.FirstName,
				newApp.LastName,
				newApp.ResumeFile,
				newApp.Phone,
				newApp.Gender,
				newApp.Github,
				newApp.Location,
				newApp.OtherMajor,
				newApp.IsFirstHackathon,
				newApp.WhyBM,
				newApp.ProjIdea,
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
			_, err = tx.Exec(`UPDATE bm_applications
			SET
				school = $1,
				major = $2,
				graduation_year = $3,
				first_name = $4,
				last_name = $5,
				resume_file = $6,
				phone = $7,
				gender = $8,
				github = $9,
				location = $10,
				other_major = $11,
				is_first_hackathon = $12,
				why_bm = $13,
				proj_idea = $14,
				tac_18_or_older = $15,
				tac_mlh_code_of_conduct = $16,
				tac_mlh_contest_and_privacy = $17
			WHERE user_id = $18`,
				newApp.School,
				newApp.Major,
				newApp.GraduationYear,
				newApp.FirstName,
				newApp.LastName,
				newApp.ResumeFile,
				newApp.Phone,
				newApp.Gender,
				newApp.Github,
				newApp.Location,
				newApp.OtherMajor,
				newApp.IsFirstHackathon,
				newApp.WhyBM,
				newApp.Is18OrOlder,
				newApp.MLHCodeOfConduct,
				newApp.MLHContestAndPrivacy,
				newApp.UserID,
			)
		} else {
			_, err = tx.Exec(`UPDATE bm_applications
			SET
				school = $1,
				major = $2,
				graduation_year = $3,
				first_name = $4,
				last_name = $5,
				phone = $6,
				gender = $7,
				github = $8,
				location = $9,
				other_major = $10,
				is_first_hackathon = $11,
				why_bm = $12,
				proj_idea= $13,
				tac_18_or_older = $14,
				tac_mlh_code_of_conduct = $15,
				tac_mlh_contest_and_privacy = $16
			WHERE user_id = $17`,
				newApp.School,
				newApp.Major,
				newApp.GraduationYear,
				newApp.FirstName,
				newApp.LastName,
				newApp.Phone,
				newApp.Gender,
				newApp.Github,
				newApp.Location,
				newApp.OtherMajor,
				newApp.IsFirstHackathon,
				newApp.WhyBM,
				newApp.ProjIdea,
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
			accepted_at,
			rsvp,
			school,
			major,
			graduation_year,
			first_name,
			last_name,
			resume_file,
			phone,
			gender,
			github,
			location,
			other_major,
			is_first_hackathon,
			why_bm,
			proj_idea,
			tac_18_or_older,
			tac_mlh_code_of_conduct,
			tac_mlh_contest_and_privacy
		FROM bm_applications
		WHERE user_id = $1`, uid).Scan(
		&dba.ID,
		&dba.UserID,
		&dba.Decision,
		&dba.AcceptedAt,
		&dba.RSVP,
		&dba.School,
		&dba.Major,
		&dba.GraduationYear,
		&dba.FirstName,
		&dba.LastName,
		&dba.ResumeFile,
		&dba.Phone,
		&dba.Gender,
		&dba.Github,
		&dba.Location,
		&dba.OtherMajor,
		&dba.IsFirstHackathon,
		&dba.WhyBM,
		&dba.ProjIdea,
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

// returns number of applications in the database
func (s *ApplicationService) GetApplicationCount() int {
	var count int

	tx, err := s.DB.Begin()
	if err != nil {
		return -1
	}

	row := tx.QueryRow("SELECT COUNT(*) FROM bm_applications")
	err = row.Scan(&count)

	if err != nil {
		tx.Rollback()
		return -1
	}

	err = tx.Commit()

	// indicates commiting transaction failed
	if err != nil {
		return -1
	}
	return count
}
