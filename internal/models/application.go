package models

import (
	"mime/multipart"
	"net/http"
	"time"
)

// Validation errors
var (
	ErrMissingSchool         = &ModelError{"please enter your school's name"}
	ErrMissingGender         = &ModelError{"please enter your gender"}
	ErrMissingMajor          = &ModelError{"please enter your major"}
	ErrMissingGraduationYear = &ModelError{"please enter your graduation year"}
	ErrMissingRace           = &ModelError{"please enter your race"}
	ErrMissingShirtSize      = &ModelError{"please enter your shirt size"}
	ErrMissingTACAgree       = &ModelError{"please agree to the terms and conditions"}

	// Validation errors when form paring
	ErrMissingResume  = &ModelError{"please upload a resume"}
	ErrResumeTooLarge = &ModelError{"resume upload is too large"}
)

const (
	DecisionAwaiting = iota
	DecisionRejected
	DecisionWaitlist
	DecisionAccepted
)

const (
	// 16MiB
	maxResumeSize = 16 << 20
)

// An Application is an application.  What do you want from me?
type Application struct {
	ID                   int
	Decision             int
	UserID               int
	School               string
	Gender               string
	Major                string
	GraduationYear       string
	DietaryRestrictions  string
	Github               string
	Linkedin             string
	ResumeFile           string
	Resume               *multipart.FileHeader // Stored in S3, not db
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

// Validate checks if an Application has all the necessary fields. Validation
// of resume uploads happens in application_service.go.
func (a *Application) Validate() error {
	if a.School == "" {
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
	a.School = r.FormValue("school")
	a.Gender = r.FormValue("gender")
	a.Major = r.FormValue("major")
	a.GraduationYear = r.FormValue("graduation-year")
	a.DietaryRestrictions = r.FormValue("dietary-restrictions")
	a.Github = r.FormValue("github")
	a.Linkedin = r.FormValue("linkedin")

	// If no file was uploaded then set ResumeFile to empty string and let
	// application_service decide what to do.  If there's already a ResumeFile
	// in the db, then they've already uploaded a resume but just haven't updated
	// it with this post request.
	_, header, err := r.FormFile("resume")
	if err != nil {
		if err != http.ErrMissingFile {
			return err
		}
	} else {
		// New file was uploaded
		a.ResumeFile = header.Filename
		a.Resume = header

		// Make sure size is reasonable
		if a.Resume.Size > maxResumeSize {
			return ErrResumeTooLarge
		}
	}

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
