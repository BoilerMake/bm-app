package models

import "errors"

var (
	ErrUserNotFound  = errors.New("user was not found")
	ErrEmailInUse    = errors.New("email is already in use")
	ErrRequiredField = errors.New("required field is missing")
)

const (
	RoleHacker = iota
	RoleSponsor
	RoleExec
)

// A User is an account stored in the database.
type User struct {
	ID           string `json:"id"`        // NOT NULL
	Role         int    `json:"role"`      // NOT NULL
	Email        string `json:"email"`     // NOT NULL
	PasswordHash string `json:"-"`         // NOT NULL
	FirstName    string `json:"firstName"` // NOT NULL
	LastName     string `json:"lastName"`  // NOT NULL
	Phone        string `json:"phone"`

	ProjectIdea string   `json:"projectIdea"`
	TeamMembers []string `json:"teamMembers"`
}

// A UserService defines an interface between the user struct (AKA the model)
// and its representation in our database.  Abstracting it to an interface
// makes it database independent, which helps with testing.
type UserService interface {
	Create(u *User) error
	GetById(id string) (*User, error)
	GetAll() (*[]User, error)
	Update(u *User) error
}
