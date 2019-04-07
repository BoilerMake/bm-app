package models

import (
	"errors"
	"time"

	"github.com/BoilerMake/new-backend/pkg/argon2"

	jwt "github.com/dgrijalva/jwt-go"
)

// Authentication errors
var (
	ErrUserNotFound   = errors.New("user was not found")
	ErrEmailInUse     = errors.New("email is already in use")
	ErrRequiredField  = errors.New("required field is missing")
	ErrIncorrectLogin = errors.New("email or password is incorrect")
)

// Validation errors
var (
	ErrEmptyEmail           = errors.New("email is empty")
	ErrEmptyPassword        = errors.New("password is empty")
	ErrEmptyPasswordConfirm = errors.New("password confirmation is empty")
	ErrEmptyFirstName       = errors.New("first name is empty")
	ErrEmptyLastName        = errors.New("last name is empty")
	ErrPasswordConfirm      = errors.New("password and confirmation password do not match")
)

const (
	RoleHacker = iota
	RoleSponsor
	RoleExec
)

// TODO this kinda goes with config as mentioned in cmd/server/main.go.  Having
// Package globals like this just doesn't sit right with me right now.
var (
	JWTSigningKey []byte
	JWTIssuer     string
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

	Phone string `json:"phone"`

	ProjectIdea string   `json:"projectIdea"`
	TeamMembers []string `json:"teamMembers"`
}

// GetJWT creates a JWT from a User.
func (u *User) GetJWT() (tokenString string, err error) {
	// TODO make sure our claims are up to par, check out:
	// https://tools.ietf.org/html/rfc7519#section-4.1
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   JWTIssuer,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"id":    u.ID,
		"role":  u.Role,
		"email": u.Email,
	})

	tokenString, err = token.SignedString(JWTSigningKey)
	if err != nil {
		return tokenString, err
	}

	return tokenString, err
}

// Validate checks if a User has all the necessary fields.
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrEmptyEmail
	} else if u.Password == "" {
		return ErrEmptyPassword
	} else if u.PasswordConfirm == "" {
		return ErrEmptyPasswordConfirm
	} else if u.PasswordConfirm != u.Password {
		return ErrPasswordConfirm
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
	Signup(u *User) (int, error)
	Login(u *User) error
	GetById(id string) (*User, error)
	GetByEmail(id string) (*User, error)
	GetAll() (*[]User, error)
	Update(u *User) error
	ResetPassword(email string) error
}
