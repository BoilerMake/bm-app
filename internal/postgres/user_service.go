// TODO If there's no rows altered does it return an err?
// TODO When using them right, transactions are usually a better way to deal
// with actions in the database.  You get rollback and some other cool stuff.
// Something to consider if you're looking for something to do.
package postgres

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"github.com/BoilerMake/new-backend/internal/models"
	"github.com/BoilerMake/new-backend/pkg/argon2"

	"github.com/lib/pq"
)

const passwordResetTokenLength int = 32
const userIDTokenLength int = 5

// Time in minutes
const tokenExpiryTime int = 15

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

	IsActive         sql.NullBool
	ConfirmationCode sql.NullString
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

		IsActive:         u.IsActive.Bool,
		ConfirmationCode: u.ConfirmationCode.String,
	}
}

// Signup inserts a new user into the database.
// TODO test cases:
// - Empty fields
// - Same email
// - nil user
func (s *UserService) Signup(u *models.User) (id int, code string, err error) {
	// Generate confirmation code

	code, err = GenerateRandomString(32)
	if err != nil {
		return id, code, err
	}

	err = u.Validate()
	if err != nil {
		return id, code, err
	}

	err = u.HashPassword()
	if err != nil {
		return id, code, err
	}

	err = s.DB.QueryRow(`INSERT INTO users (
			role,
			email,
			password_hash,
			first_name,
			last_name,
			phone,
			is_active,
			confirmation_code
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;`,
		u.Role,
		u.Email,
		u.PasswordHash,
		u.FirstName,
		u.LastName,
		u.Phone,
		false,
		code,
	).Scan(&id)

	// Check postgres specific error
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "unique_violation":
			return id, code, models.ErrEmailInUse
		case "not_null_violation":
			return id, code, models.ErrRequiredField
		default:
			return id, code, pgerr
		}
	}

	return id, code, err
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
}

// GetById returns a single user with the given id.
func (s *UserService) GetById(id int) (*models.User, error) {
	var dbu dbUser

	err := s.DB.QueryRow(`SELECT
		id,
		role,
		email,
		password_hash,
		first_name,
		last_name,
		phone,
		is_active,
		confirmation_code
	FROM users
	WHERE id = $1`, id).Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.PasswordHash, &dbu.first_name, &dbu.last_name, &dbu.Phone, &dbu.IsActive, &dbu.ConfirmationCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrIncorrectLogin
		} else {
			return nil, err
		}
	}

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
		is_active,
		confirmation_code
	FROM users
	WHERE email = $1`, email).Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.PasswordHash, &dbu.first_name, &dbu.last_name, &dbu.Phone, &dbu.IsActive, &dbu.ConfirmationCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrIncorrectLogin
		} else {
			return nil, err
		}
	}

	// TODO if there's an err dbu will likely be nil so toModel will panic.
	// Seems like toModel needs to check for nil and maybe return an err.
	// Definitely something to test.
	return dbu.toModel(), err
}

// GetByCode returns a single user with the given email confirmation code
func (s *UserService) GetByCode(code string) (*models.User, error) {
	var dbu dbUser

	err := s.DB.QueryRow(`SELECT
		id,
		role,
		email,
		password_hash,
		first_name,
		last_name,
		phone,
		is_active,
		confirmation_code
	FROM users
	WHERE confirmation_code = $1`, code).Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.PasswordHash, &dbu.first_name, &dbu.last_name, &dbu.Phone, &dbu.IsActive, &dbu.ConfirmationCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrInvalidConfirmationCode
		} else {
			return nil, err
		}
	}

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
		is_active,
		confirmation_code
	FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbu dbUser
		err = rows.Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.first_name, &dbu.last_name, &dbu.Phone, &dbu.IsActive, &dbu.ConfirmationCode)
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
func (s *UserService) Update(u *models.User) error {
	_, err := s.DB.Exec(`UPDATE users
	SET
		role = $1,
		email = $2,
		password_hash = $3,
		first_name = $4,
		last_name = $5,
		phone = $6,
		is_active = $7,
		confirmation_code = $8
	WHERE id = $9`, u.Role, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.Phone, u.IsActive, u.ConfirmationCode, u.ID)

	if err == sql.ErrNoRows {
		return models.ErrIncorrectLogin
	}
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

// Creates the user's token for a password reset
func (s *UserService) GetPasswordReset(email string) (string, error) {
	if email == "" {
		return "", models.ErrEmptyEmail
	}

	token, err := GenerateRandomString(passwordResetTokenLength)
	if err != nil {
		return "", err
	}
	hashedToken, err := argon2.DefaultParameters.HashPassword(token)
	if err != nil {
		return "", err
	}
	tokenID, err := GenerateRandomString(userIDTokenLength)
	if err != nil {
		return "", err
	}

	_, err = s.DB.Exec(`
	INSERT INTO
		password_reset_tokens (uid, tokenID, hashedToken)
	VALUES
		((SELECT id FROM users WHERE email = $1), $2, $3);`, email, tokenID, hashedToken)

	// User should not know if the email exists
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "not_null_violation":
			return "", nil
		}
	}

	userToken := tokenID + token

	return userToken, nil
}

// ResetPassword resets the user's password
func (s *UserService) ResetPassword(token string, password string) error {
	if len(token) < userIDTokenLength+passwordResetTokenLength {
		return models.ErrInvalidToken
	}
	tokenID := token[:userIDTokenLength]
	userToken := token[userIDTokenLength:]

	id := -1
	var uid int
	var hashedToken string
	var createdAt time.Time
	var now time.Time
	// TODO check all rows with the same tokenID
	// Multiple people can have the same tokenID
	row := s.DB.QueryRow(`SELECT id, uid, hashedToken, created_at, current_timestamp FROM password_reset_tokens WHERE tokenID = $1`, tokenID)
	err := row.Scan(&id, &uid, &hashedToken, &createdAt, &now)
	elapsed := now.Sub(createdAt).Minutes()

	// Check if token is expired
	if elapsed > float64(tokenExpiryTime) || elapsed < 0 {
		return models.ErrExpiredToken
	}

	// Not sure what error to return
	// User should not know if the email exists
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		// User not in db (log internally)
		// No error should be returned to user
		case "not_null_violation":
			return nil
		default:
			return pgerr
		}
	}

	if argon2.CheckPassword(userToken, hashedToken) {
		s.TokenChangePassword(id, uid, password)
		return nil
	}

	return models.ErrInvalidToken
}

// TokenChangePassword changes password then deletes the token used
func (s *UserService) TokenChangePassword(id int, uid int, password string) error {
	passwordHash, err := argon2.DefaultParameters.HashPassword(password)
	if err != nil {
		return err
	}

	// Remove any trace of unhashed password (can never be too safe)
	password = ""

	_, err = s.DB.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`, passwordHash, uid)
	_, err = s.DB.Exec(`DELETE FROM password_reset_tokens WHERE id=$1`, id)

	return err
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
