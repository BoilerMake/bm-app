// TODO consider sqlx: https://github.com/jmoiron/sqlx/
// ^ Should make most methods a lot prettier
// TODO if there's no rows altered does it return an err?
package postgres

import (
	"database/sql"
	"fmt"

	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/lib/pq"
)

// NOTE this is where postgres error codes come from:
// https://github.com/lib/pq/blob/9eb73efc1fcc404148b56765b0d3f61d9a5ef8ee/error.go#L78

// A dbUser is like a models.User, but it can read in null fields from the
// database without panicking.
type dbUser struct {
	// TODO do we still need sql.Null* types for NOT NULL columns? I'm guessing no
	ID           sql.NullString
	Role         sql.NullInt64
	Email        sql.NullString
	PasswordHash sql.NullString
	FirstName    sql.NullString
	LastName     sql.NullString
	Phone        sql.NullString

	ProjectIdea sql.NullString
	// TODO figure out how to read in array from db
	TeamMember1 sql.NullString
	TeamMember2 sql.NullString
	TeamMember3 sql.NullString
}

// toModel converts a database specific dbUser to the more generic User struct.
func (u *dbUser) toModel() *models.User {
	// TODO test to see what happens if u is nil
	return &models.User{
		ID:           u.ID.String,
		Role:         int(u.Role.Int64),
		Email:        u.Email.String,
		PasswordHash: u.PasswordHash.String,
		FirstName:    u.FirstName.String,
		LastName:     u.LastName.String,
		Phone:        u.Phone.String,

		ProjectIdea: u.ProjectIdea.String,
		TeamMember1: u.TeamMember1.String,
		TeamMember2: u.TeamMember2.String,
		TeamMember3: u.TeamMember3.String,
	}
}

// UserService is a PostgreSQL implementation of models.UserService
type UserService struct {
	DB *sql.DB
}

// Create inserts a new user into the database.
// TODO test cases:
// - Empty fields
// - Same email
// - nil user
func (s *UserService) Create(u *models.User) error {
	// TODO convert teammembern into array
	_, err := s.DB.Exec(`INSERT INTO users(
			role, 
			email, 
			password_hash, 
			firstname, 
			lastname, 
			phone, 
			projectidea, 
			teammember1, 
			teammember2, 
			teammember3) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, u.Role, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Phone, u.ProjectIdea, u.TeamMember1, u.TeamMember2, u.TeamMember3)

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
		teammember1, 
		teammember2, 
		teammember3 
	FROM users 
	WHERE id = $1`, id).Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.PasswordHash, &dbu.FirstName, &dbu.LastName, &dbu.Phone, &dbu.ProjectIdea, &dbu.TeamMember1, &dbu.TeamMember2, &dbu.TeamMember3)

	// FIXME i think if there's an err dbu will be nil so toModel will panic.
	// Seems like toModel needs to check for nil and maybe return an err.
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
		teammember1, 
		teammember2, 
		teammember3 
	FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbu dbUser
		err = rows.Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.FirstName, &dbu.LastName, &dbu.Phone, &dbu.ProjectIdea, &dbu.TeamMember1, &dbu.TeamMember2, &dbu.TeamMember3)
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
// TODO being able to update email has some implications to it:
// - NOT NULL/UNIQUE
// - Reconfirming their new email
// - Removing old email from mailing lists and adding new one
//   - Do we even have mailing lists for this to be a problem with?
// - Obvious answer seems to just not allow changing email (what we've always done)
func (s *UserService) Update(u *models.User) error {
	_, err := s.DB.Exec(`UPDATE users
	SET 
		role = $1, 
		email = $2, 
		password_hash = $3, 
		firstname = $3, 
		lastname = $5, 
		phone = $6, 
		projectidea = $7, 
		teammember1 = $8, 
		teammember2 = $9, 
		teammember3 = $10
	WHERE id = $11`, u.Role, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Phone, u.ProjectIdea, u.TeamMember1, u.TeamMember2, u.TeamMember3, u.ID)

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
