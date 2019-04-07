// TODO If there's no rows altered does it return an err?
// TODO When using them right, transactions are usually a better way to deal
// with actions in the database.  You get rollback and some other cool stuff.
// Something to consider if you're looking for something to do.
package postgres

import (
	"database/sql"
	"fmt"

	"github.com/BoilerMake/new-backend/internal/models"

	"github.com/lib/pq"
)

// UserService is a PostgreSQL implementation of models.UserService
type UserService struct {
	DB *sql.DB
}

// A dbUser is like a models.User, but it can read in null fields from the
// database without panicking.
type dbUser struct {
	// TODO do we still need sql.Null* types for NOT NULL columns? I would guess no
	ID           sql.NullInt64
	Role         sql.NullInt64
	Email        sql.NullString
	PasswordHash sql.NullString
	first_name   sql.NullString
	last_name    sql.NullString
	Phone        sql.NullString

	ProjectIdea sql.NullString
	TeamMembers []string
}

// toModel converts a database specific dbUser to the more generic User struct.
func (u *dbUser) toModel() *models.User {
	// TODO test to see what happens if u is nil
	return &models.User{
		ID:           int(u.ID.Int64),
		Role:         int(u.Role.Int64),
		Email:        u.Email.String,
		PasswordHash: u.PasswordHash.String,
		FirstName:    u.first_name.String,
		LastName:     u.last_name.String,
		Phone:        u.Phone.String,

		ProjectIdea: u.ProjectIdea.String,
		TeamMembers: u.TeamMembers,
	}
}

// Signup inserts a new user into the database.
// TODO test cases:
// - Empty fields
// - Same email
// - nil user
func (s *UserService) Signup(u *models.User) (id int, err error) {
	err = u.Validate()
	if err != nil {
		return id, err
	}

	err = u.HashPassword()
	if err != nil {
		return id, err
	}

	err = s.DB.QueryRow(`INSERT INTO users (
			role, 
			email, 
			password_hash, 
			first_name, 
			last_name, 
			phone, 
			project_idea, 
			team_members
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id;`,
		u.Role,
		u.Email,
		u.PasswordHash,
		u.FirstName,
		u.LastName,
		u.Phone,
		u.ProjectIdea,
		pq.Array(u.TeamMembers),
	).Scan(&id)

	// Check postgres specific error
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "unique_violation":
			return id, models.ErrEmailInUse
		case "not_null_violation":
			return id, models.ErrRequiredField
		default:
			return id, pgerr
		}
	}

	return id, err
}

// Login checks to see if a passed user's password matches the one recorded in
// the database.  It expects to receive a user object with at least an email
// and a password.  When Login fails it will return an error,  otherwise it
// will return nil.
func (s *UserService) Login(u *models.User) error {
	dbu, err := s.GetByEmail(u.Email)
	if err != nil {
		return err
	}

	if dbu.CheckPassword(u.Password) {
		// We know the login was valid now, but the user constructed from the
		// request doesn't have all its fields filled in.  It will need those
		// fields to create a new JWT, so we'll set them here.
		*u = *dbu
		return nil
	} else {
		return models.ErrIncorrectLogin
	}

	// UNREACHABLE CODE
	//return err
}

// GetById returns a single user with the given id.
func (s *UserService) GetById(id string) (*models.User, error) {
	var dbu dbUser

	err := s.DB.QueryRow(`SELECT 
		id, 
		role, 
		email, 
		password_hash, 
		first_name, 
		last_name, 
		phone, 
		project_idea, 
		team_members
	FROM users 
	WHERE id = $1`, id).Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.PasswordHash, &dbu.first_name, &dbu.last_name, &dbu.Phone, &dbu.ProjectIdea, pq.Array(&dbu.TeamMembers))

	// TODO if there's an err dbu will likely be nil so toModel will panic.
	// Seems like toModel needs to check for nil and maybe return an err.
	// Definitely something to test.
	return dbu.toModel(), err
}

// GetByEmail returns a single user with the given email.
func (s *UserService) GetByEmail(email string) (*models.User, error) {
	var dbu dbUser

	err := s.DB.QueryRow(`SELECT 
		id, 
		role, 
		email, 
		password_hash, 
		first_name, 
		last_name, 
		phone, 
		project_idea, 
		team_members
	FROM users 
	WHERE email = $1`, email).Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.PasswordHash, &dbu.first_name, &dbu.last_name, &dbu.Phone, &dbu.ProjectIdea, pq.Array(&dbu.TeamMembers))

	if err != nil {
		return nil, err
	}

	// TODO if there's an err dbu will likely be nil so toModel will panic.
	// Seems like toModel needs to check for nil and maybe return an err.
	// Definitely something to test.
	return dbu.toModel(), err
}

// GetAll finds and returns all user in the database.
func (s *UserService) GetAll() (u *[]models.User, err error) {
	rows, err := s.DB.Query(`SELECT 
		id, 
		role, 
		email, 
		first_name, 
		last_name, 
		phone, 
		project_idea, 
		team_members 
	FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbu dbUser
		err = rows.Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.first_name, &dbu.last_name, &dbu.Phone, &dbu.ProjectIdea, &dbu.TeamMembers)
		if err != nil {
			return nil, fmt.Errorf("failed to get all users: %v", err)
		}

		*u = append(*u, *dbu.toModel())
	}

	return u, err
}

// Update changes the values of a user in the database.  It will update all
// fields in the user model given to it, including nil or zero fields.  Don't
// give it a user with only the changes you want to make.
//
// TODO being able to update email has some implications to it:
// - NOT NULL/UNIQUE
// - Reconfirming their new email
// - Removing old email from mailing lists(?) and adding new one
//   - Do we even have mailing lists for this to be a problem with?
// - Obvious answer seems to just not allow changing email (what we've always done)
func (s *UserService) Update(u *models.User) error {
	_, err := s.DB.Exec(`UPDATE users
	SET 
		role = $1, 
		email = $2, 
		password_hash = $3, 
		first_name = $3, 
		last_name = $5, 
		phone = $6, 
		project_idea = $7, 
		team_members = $8, 
	WHERE id = $11`, u.Role, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Phone, u.ProjectIdea, pq.Array(u.TeamMembers), u.ID)

	// TODO make sure when an exec fails it doesn't have an effect on the db
	// Check postgres specific error
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "unique_violation":
			return models.ErrEmailInUse
		case "not_null_violation":
			return models.ErrRequiredField
		default:
			return pgerr
		}
	}

	return err
}

// TODO remove token automatically if expired
// Finding email happends in this method, could be neater
func (s *UserService) ResetPassword(email string) error {

	randomToken := "RandomToken"
	_, err := s.DB.Exec(`
	INSERT INTO
		password_reset_tokens (uid, token, valid_until)
	VALUES
		((SELECT id FROM users WHERE email = $1), $2, current_timestamp + interval '1 hour')
	ON CONFLICT (uid)
	DO UPDATE
	SET
	 	uid = (SELECT id FROM users WHERE email = $1), 
    	token = $2, 
		valid_until = current_timestamp + interval '1 hour';`, email, randomToken)

	// Not sure what to return
	// User should not know if the email exists
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "not_null_violation":
			return models.ErrRequiredField
		default:
			return pgerr
		}
	}

	return nil
}
