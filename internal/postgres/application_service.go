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
	ID                  sql.NullInt64
	UserID              sql.NullInt64
	Decision            sql.NullInt64
	EmailedDecision     sql.NullBool
	CheckedInAt         pq.NullTime
	RSVP                sql.NullBool
	School              sql.NullString
	Gender              sql.NullString
	Major               sql.NullString
	GraduationYear      sql.NullString
	DietaryRestrictions sql.NullString
	Github              sql.NullString
	Linkedin            sql.NullString
	// TODO has resume
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
		ID:                  int(a.ID.Int64),
		UserID:              int(a.UserID.Int64),
		Decision:            int(a.Decision.Int64),
		EmailedDecision:     a.EmailedDecision.Bool,
		CheckedInAt:         a.CheckedInAt.Time,
		RSVP:                a.RSVP.Bool,
		School:              a.School.String,
		Gender:              a.Gender.String,
		Major:               a.Major.String,
		GraduationYear:      a.GraduationYear.String,
		DietaryRestrictions: a.DietaryRestrictions.String,
		Github:              a.Github.String,
		Linkedin:            a.Linkedin.String,
		// TODO has_resume
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
func (s *ApplicationService) CreateOrUpdate(a *models.Application) (err error) {
	err = a.Validate()
	if err != nil {
		return err
	}

	// TODO has_resume
	// If there's a conflict on user_id (a row already exists with that user_id),
	// then we know that user already has an application so we can just update that
	_, err = s.DB.Exec(`INSERT INTO bm7_applications (
			user_id,
			school,
			gender,
			major,
			graduation_year,
			dietary_restrictions,
			github,
			linkedin,
			is_first_hackathon,
			race,
			shirt_size,
			project_idea,
			team_members,
			tac_18_or_older,
			tac_mlh_code_of_conduct,
			tac_mlh_contest_and_privacy
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (user_id) DO UPDATE SET
			user_id = $1,
			school = $2,
			gender = $3,
			major = $4,
			graduation_year = $5,
			dietary_restrictions = $6,
			github = $7,
			linkedin = $8,
			is_first_hackathon = $9,
			race = $10,
			shirt_size = $11,
			project_idea = $12,
			team_members = $13,
			tac_18_or_older = $14,
			tac_mlh_code_of_conduct = $15,
			tac_mlh_contest_and_privacy = $16;`,
		a.UserID,
		a.School,
		a.Gender,
		a.Major,
		a.GraduationYear,
		a.DietaryRestrictions,
		a.Github,
		a.Linkedin,
		a.IsFirstHackathon,
		a.Race,
		a.ShirtSize,
		a.ProjectIdea,
		pq.Array(a.TeamMembers),
		a.Is18OrOlder,
		a.MLHCodeOfConduct,
		a.MLHContestAndPrivacy,
	)

	return err
}

// GetByUserID returns a single Application with the given user id.
func (s *ApplicationService) GetByUserID(uid int) (*models.Application, error) {
	var dba dbApplication

	// TODO has_resume
	err := s.DB.QueryRow(`SELECT
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
		return nil, err
	}

	// TODO if there's an err dbu will likely be nil so toModel will panic.
	// Seems like toModel needs to check for nil and maybe return an err.
	// Definitely something to test.
	return dba.toModel(), err
}
