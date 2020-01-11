package models

import (
	"net/http"
	"strings"

	"github.com/BoilerMake/bm-app/pkg/argon2"
	"github.com/BoilerMake/bm-app/pkg/flash"

	"github.com/gorilla/sessions"
)

// Authentication errors
var (
	ErrUserNotFound            = &ModelError{"User not found.", flash.Info}
	ErrEmailInUse              = &ModelError{"Email is already in use.", flash.Info}
	ErrRequiredField           = &ModelError{"Required field is missing.", flash.Info}
	ErrIncorrectLogin          = &ModelError{"Email or password is incorrect.", flash.Info}
	ErrNotLoggedIn             = &ModelError{"Please create an account or log in first.", flash.Info}
	ErrAlreadyLoggedIn         = &ModelError{"You're already logged in.", flash.Info}
	ErrInvalidConfirmationCode = &ModelError{"Email confirmation code no longer valid.", flash.Info}
)

// Validation errors
var (
	ErrEmptyEmail           = &ModelError{"Email is empty", flash.Info}
	ErrInvalidEmail         = &ModelError{"Email is invalid", flash.Info}
	ErrEmailNotFound        = &ModelError{"We don't recognize that email.", flash.Info}
	ErrEmptyPassword        = &ModelError{"Password is empty", flash.Info}
	ErrEmptyPasswordConfirm = &ModelError{"Password confirmation is empty", flash.Info}
	ErrEmptyFirstName       = &ModelError{"First name is empty", flash.Info}
	ErrEmptyLastName        = &ModelError{"Last name is empty", flash.Info}
	ErrPasswordConfirm      = &ModelError{"Password and confirmation password do not match", flash.Info}
	ErrPasswordTooShort     = &ModelError{"Password must be at least 3 characters", flash.Info}
)

// Password Reset errors
var (
	ErrInvalidToken = &ModelError{"Password reset token is invalid", flash.Info}
	ErrExpiredToken = &ModelError{"Password reset token has expired", flash.Info}
)

const (
	RoleHacker = iota
	RoleSponsor
	RoleExec
)

// A User is an account stored in the database.
type User struct {
	ID   int `json:"id"`   // NOT NULL
	Role int `json:"role"` // NOT NULL

	Email string `json:"email"` // NOT NULL

	// Password and PasswordConfirm should only ever be in memory, never in the db
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	PasswordHash    string `json:"-"` // NOT NULL

	FirstName string `json:"firstName"` // NOT NULL
	LastName  string `json:"lastName"`  // NOT NULL

	IsActive         bool   `json:"isActive"`
	ConfirmationCode string `json:"confirmationCode"`
}

// EmailModel struct for password reset emails
type EmailModel struct {
	Email string `json:"email"`
}

// PasswordResetPayload struct for resetting passwords
type PasswordResetPayload struct {
	UserToken   string `json:"token"`
	NewPassword string `json:"newPassword"`
}

// SetSession sets all the session values for a user
func (u *User) SetSession(session *sessions.Session) {
	session.Values["ID"] = u.ID
	session.Values["EMAIL"] = u.Email
	session.Values["ROLE"] = u.Role
}

// FromFormData converts a user from a request's FormData to a models.User
// struct.
func (u *User) FromFormData(r *http.Request) {
	u.Email = r.FormValue("email")
	u.Email = strings.ToLower(u.Email)

	u.Password = r.FormValue("password")
	u.PasswordConfirm = r.FormValue("password-confirm")

	u.FirstName = r.FormValue("first-name")
	u.LastName = r.FormValue("last-name")
}

// Validate checks if a User has all the necessary fields.
func (u *User) Validate() error {
	u.Email = strings.ToLower(u.Email)

	if u.Email == "" {
		return ErrEmptyEmail
	} else if u.Password == "" {
		return ErrEmptyPassword
	} else if u.PasswordConfirm == "" {
		return ErrEmptyPasswordConfirm
	} else if u.PasswordConfirm != u.Password {
		return ErrPasswordConfirm
	} else if len(u.Password) < 3 {
		return ErrPasswordTooShort
	} else if u.FirstName == "" {
		return ErrEmptyFirstName
	} else if u.LastName == "" {
		return ErrEmptyLastName
	}

	return nil
}

// HashPassword hashes a User's password and sets its PasswordHash field.  It
// also empties its Password and PasswordConfirm fields.
func (u *User) HashPassword() error {
	passwordHash, err := argon2.DefaultParameters.HashPassword(u.Password)
	if err != nil {
		return err
	}

	// Remove any trace of unhashed password (can never be too safe)
	u.Password = ""
	u.PasswordConfirm = ""

	u.PasswordHash = passwordHash
	return nil
}

// CheckPassword compares a User's hashed password to a string.
func (u *User) CheckPassword(password string) bool {
	return argon2.CheckPassword(password, u.PasswordHash)
}

// A UserService defines an interface between the user struct (AKA the model)
// and its representation in our database.  Abstracting it to an interface
// makes it database independent, which helps with testing.
type UserService interface {
	Signup(u *User) (int, string, error)
	Login(u *User) error
	GetById(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByCode(code string) (*User, error)
	Update(u *User) error
	GetPasswordReset(email string) (string, error)
	ResetPassword(token string, newPassword string) error
	GetUserCount() int
}
