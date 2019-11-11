package postgres

import (
	"crypto/rand"
	"database/sql"
	"strings"
	"time"

	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/argon2"

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
	ID           sql.NullInt64
	Role         sql.NullInt64
	Email        sql.NullString
	PasswordHash sql.NullString
	FirstName    sql.NullString
	LastName     sql.NullString

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
		FirstName:    u.FirstName.String,
		LastName:     u.LastName.String,

		IsActive:         u.IsActive.Bool,
		ConfirmationCode: u.ConfirmationCode.String,
	}
}

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

	tx, err := s.DB.Begin()
	if err != nil {
		return id, code, err
	}

	err = tx.QueryRow(`INSERT INTO users (
			role,
			email,
			password_hash,
			first_name,
			last_name,
			is_active,
			confirmation_code
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`,
		u.Role,
		u.Email,
		u.PasswordHash,
		u.FirstName,
		u.LastName,
		false,
		code,
	).Scan(&id)

	// Check postgres specific error
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "unique_violation":
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return id, code, rollbackErr
			}
			return id, code, models.ErrEmailInUse
		case "not_null_violation":
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return id, code, rollbackErr
			}
			return id, code, models.ErrRequiredField
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return id, code, rollbackErr
			}
			return id, code, pgerr
		}
	}

	// Switch statement above only checks for postgres specific errors, so here we'll check
	// for a generic error
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return id, code, rollbackErr
		}
		return id, code, err
	}

	err = tx.Commit()
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
func (s *UserService) GetById(id int) (u *models.User, err error) {
	var dbu dbUser

	tx, err := s.DB.Begin()
	if err != nil {
		return u, err
	}

	err = tx.QueryRow(`SELECT
		id,
		role,
		email,
		password_hash,
		first_name,
		last_name,
		is_active,
		confirmation_code
	FROM users
	WHERE id = $1`, id).Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.PasswordHash, &dbu.FirstName, &dbu.LastName, &dbu.IsActive, &dbu.ConfirmationCode)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}

		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		} else {
			return nil, err
		}
	}

	err = tx.Commit()
	return dbu.toModel(), err
}

// GetByEmail returns a single user with the given email.
func (s *UserService) GetByEmail(email string) (u *models.User, err error) {
	email = strings.ToLower(email)
	var dbu dbUser

	tx, err := s.DB.Begin()
	if err != nil {
		return u, err
	}

	err = tx.QueryRow(`SELECT
		id,
		role,
		email,
		password_hash,
		first_name,
		last_name,
		is_active,
		confirmation_code
	FROM users
	WHERE email = $1`, email).Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.PasswordHash, &dbu.FirstName, &dbu.LastName, &dbu.IsActive, &dbu.ConfirmationCode)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return u, rollbackErr
		}

		if err == sql.ErrNoRows {
			return nil, models.ErrEmailNotFound
		} else {
			return nil, err
		}
	}

	err = tx.Commit()
	return dbu.toModel(), err
}

// GetByCode returns a single user with the given email confirmation code
func (s *UserService) GetByCode(code string) (u *models.User, err error) {
	var dbu dbUser

	tx, err := s.DB.Begin()
	if err != nil {
		return u, err
	}

	err = tx.QueryRow(`SELECT
		id,
		role,
		email,
		password_hash,
		first_name,
		last_name,
		is_active,
		confirmation_code
	FROM users
	WHERE confirmation_code = $1`, code).Scan(&dbu.ID, &dbu.Role, &dbu.Email, &dbu.PasswordHash, &dbu.FirstName, &dbu.LastName, &dbu.IsActive, &dbu.ConfirmationCode)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}

		if err == sql.ErrNoRows {
			return nil, models.ErrInvalidConfirmationCode
		} else {
			return nil, err
		}
	}

	err = tx.Commit()
	return dbu.toModel(), err
}

// Update changes the values of a user in the database.  It will update all
// fields in the user model given to it, including nil or zero fields.  Don't
// give it a user with only the changes you want to make.
func (s *UserService) Update(u *models.User) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	u.Email = strings.ToLower(u.Email)

	_, err = tx.Exec(`UPDATE users
	SET
		role = $1,
		email = LOWER($2),
		password_hash = $3,
		first_name = $4,
		last_name = $5,
		is_active = $6,
		confirmation_code = $7
	WHERE id = $8`, u.Role, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.IsActive, u.ConfirmationCode, u.ID)

	// Check postgres specific error
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "unique_violation":
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return models.ErrEmailInUse
		case "not_null_violation":
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return models.ErrRequiredField
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return pgerr
		}
	}

	// Switch statement above only checks for postgres specific errors, so here we'll check
	// for a generic error
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		if err == sql.ErrNoRows {
			return models.ErrUserNotFound
		} else {
			return err
		}
	}

	err = tx.Commit()
	return err
}

// Creates the user's token for a password reset
func (s *UserService) GetPasswordReset(email string) (string, error) {
	email = strings.ToLower(email)

	if email == "" {
		return "", models.ErrEmptyEmail
	}

	dbu, err := s.GetByEmail(email)
	if err != nil {
		return "", err
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

	tx, err := s.DB.Begin()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(`INSERT INTO
		password_reset_tokens (uid, tokenID, hashedToken)
	VALUES
		($1, $2, $3);`, dbu.ID, tokenID, hashedToken)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return "", rollbackErr
		}
		return "", err
	}

	userToken := tokenID + token

	err = tx.Commit()
	return userToken, err
}

// ResetPassword resets the user's password
func (s *UserService) ResetPassword(token string, password string) error {
	if len(token) < userIDTokenLength+passwordResetTokenLength {
		return models.ErrInvalidToken
	}
	tokenID := token[:userIDTokenLength]
	userToken := token[userIDTokenLength:]

	if len(password) < 3 {
		return models.ErrPasswordTooShort
	}

	// Multiple people can have the same tokenID
	// Must check all
	rows, err := s.DB.Query("SELECT id, uid, hashedToken, created_at, current_timestamp FROM password_reset_tokens WHERE tokenID = $1", tokenID)
	if err != nil {
		return models.ErrInvalidToken
	}
	defer rows.Close()

	var id int
	var uid int
	var hashedToken string
	var createdAt time.Time
	var now time.Time
	for rows.Next() {
		err = rows.Scan(&id, &uid, &hashedToken, &createdAt, &now)
		if err != nil {
			return models.ErrInvalidToken
		}
		elapsed := now.Sub(createdAt).Minutes()

		// Check if token is expired
		if elapsed > float64(tokenExpiryTime) || elapsed < 0 {
			// The tokenID may belong to another user, so don't modify it
			continue
		}

		// Check if hash is correct
		if argon2.CheckPassword(userToken, hashedToken) {
			err = s.TokenChangePassword(id, uid, password)
			return err
		}
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

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`, passwordHash, uid)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	tx, err = s.DB.Begin()
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(`DELETE FROM password_reset_tokens WHERE id=$1`, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()
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
