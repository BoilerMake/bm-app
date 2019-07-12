package models

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

// Validation errors
var (
	ErrMissingSchool         = errors.New("please enter your school's name")
	ErrMissingGender         = errors.New("please enter your gender")
	ErrMissingMajor          = errors.New("please enter your major")
	ErrMissingGraduationYear = errors.New("please enter your graduation year")
	ErrMissingRace           = errors.New("please enter your race")
	ErrMissingShirtSize      = errors.New("please enter your shirt size")
	ErrMissingTACAgree       = errors.New("please agree to the terms and conditions")
)

// Form parsing errors
var (
	ErrUnknownSchool = errors.New("unknown school selected")
)

const (
	DecisionAwaiting = iota
	DecisionRejected
	DecisionWaitlist
	DecisionAccepted
)

// An Application is an application.  What do you want from me?
type Application struct {
	ID                   int
	Decision             int
	UserID               int
	SchoolID             int
	Gender               string
	Major                string
	GraduationYear       string
	DietaryRestrictions  string
	Github               string
	Linkedin             string
	HasResume            bool
	RSVP                 bool
	IsFirstHackathon     bool
	Race                 string
	EmailedDecision      bool
	CheckedInAt          time.Time
	ShirtSize            string
	ProjectIdea          string
	TeamMembers          []string
	Is18OrOlder          bool
	MLHCodeOfConduct     bool
	MLHContestAndPrivacy bool
}

// Validate checks if an Application has all the necessary fields.
func (a *Application) Validate() error {
	if a.SchoolID == 0 {
		return ErrMissingSchool
	} else if a.Gender == "" {
		return ErrMissingGender
	} else if a.Major == "" {
		return ErrMissingMajor
	} else if a.GraduationYear == "" {
		return ErrMissingGraduationYear
	} else if a.Race == "" {
		return ErrMissingRace
	} else if a.ShirtSize == "" {
		return ErrMissingShirtSize
	} else if !a.Is18OrOlder || !a.MLHCodeOfConduct || !a.MLHContestAndPrivacy {
		return ErrMissingTACAgree
	}

	return nil
}

// FromFormData converts an application from a request's FormData to a
// models.Application struct.
func (a *Application) FromFormData(r *http.Request) error {
	schoolID, err := strconv.Atoi(r.FormValue("school"))
	if err != nil {
		return ErrUnknownSchool
	}
	a.SchoolID = schoolID

	a.Gender = r.FormValue("gender")
	a.Major = r.FormValue("major")
	a.GraduationYear = r.FormValue("graduation-year")
	a.DietaryRestrictions = r.FormValue("dietary-restrictions")
	a.Github = r.FormValue("github")
	a.Linkedin = r.FormValue("linkedin")
	//u.HasResume = r.FormValue("has-resume") // TODO resume handling?
	a.IsFirstHackathon = r.FormValue("is-first-hackathon") == "on"
	a.Race = r.FormValue("race")
	a.ShirtSize = r.FormValue("shirt-size")
	a.ProjectIdea = r.FormValue("project-idea")
	a.TeamMembers = append(a.TeamMembers, r.FormValue("team-member-1"), r.FormValue("team-member-2"), r.FormValue("team-member-3"))

	a.Is18OrOlder = r.FormValue("is-18-or-older") == "on"
	a.MLHCodeOfConduct = r.FormValue("mlh-code-of-conduct") == "on"
	a.MLHContestAndPrivacy = r.FormValue("mlh-contest-and-privacy") == "on"

	return nil
}

// An ApplicationService defines an interface between the Application struct
// (AKA the model) and its representation in our database.  Abstracting it to
// an interface makes it database independent, which helps with testing.
type ApplicationService interface {
	CreateOrUpdate(a *Application) error
	GetByUserID(uid int) (*Application, error)
}
