package models

import (
	"errors"
	"mime/multipart"
	"net/http"
	"time"
)

// Validation errors
var (
	ErrMissingSchool         = errors.New("please enter your school's name")
	ErrMissingMajor          = errors.New("please enter your major")
	ErrMissingGraduationYear = errors.New("please enter your graduation year")
	ErrMissingGender         = errors.New("please enter your gender")
	ErrMissingRace           = errors.New("please enter your race")
	ErrMissingWhyBM          = errors.New("please enter why you want to come to BoilerMake?")
	ErrMissingTACAgree       = errors.New("please agree to the terms and conditions")

	// Validation errors when form paring
	ErrMissingResume  = errors.New("please upload a resume")
	ErrResumeTooLarge = errors.New("resume upload is too large")
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
	EmailedDecision      bool
	UserID               int
	RSVP                 bool
	CheckedInAt          time.Time
	School               string
	Major                string
	GraduationYear       string
	ResumeFile           string
	Resume               *multipart.FileHeader // Stored in S3, not db
	Gender               string
	Race                 string
	DietaryRestrictions  string
	Github               string
	Linkedin             string
	WhyBM                string
	Is18OrOlder          bool
	MLHCodeOfConduct     bool
	MLHContestAndPrivacy bool
}

// Validate checks if an Application has all the necessary fields. Validation
// of resume uploads happens in application_service.go.
func (a *Application) Validate() error {
	if a.School == "" {
		return ErrMissingSchool
	} else if a.Major == "" {
		return ErrMissingMajor
	} else if a.GraduationYear == "" {
		return ErrMissingGraduationYear
	} else if a.Gender == "" {
		return ErrMissingGender
	} else if a.Race == "" {
		return ErrMissingRace
	} else if a.WhyBM == "" {
		return ErrMissingWhyBM
	} else if !a.Is18OrOlder || !a.MLHCodeOfConduct || !a.MLHContestAndPrivacy {
		return ErrMissingTACAgree
	}

	return nil
}

// FromFormData converts an application from a request's FormData to a
// models.Application struct.
func (a *Application) FromFormData(r *http.Request) error {
	a.School = r.FormValue("school")
	a.Major = r.FormValue("major")
	a.GraduationYear = r.FormValue("graduation-year")
	a.Gender = r.FormValue("gender")
	a.Race = r.FormValue("race")
	a.DietaryRestrictions = r.FormValue("dietary-restrictions")
	a.Github = r.FormValue("github")
	a.Linkedin = r.FormValue("linkedin")
	a.WhyBM = r.FormValue("why-bm")

	// If no file was uploaded then set ResumeFile to empty string and let
	// application_service decide what to do.  If there's already a ResumeFile
	// in the db, then they've already uploaded a resume but just haven't updated
	// it with this post request.
	_, header, err := r.FormFile("resume")
	if err != nil {
		// Check if this error happened becuase request was too large
		// Kinda janky but it works
		if err.Error() == "multipart: NextPart: http: request body too large" {
			return ErrResumeTooLarge
		}

		// Otherwise only return an error if it's not because the file was missing.  Again,
		// we handle the msising resume case in the database.
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
