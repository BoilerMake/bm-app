package models

import (
	"mime/multipart"
	"net/http"
	"time"

	"github.com/BoilerMake/bm-app/pkg/flash"
)

// Validation errors
var (
	ErrMissingSchool         = &ModelError{"Please enter your school's name.", flash.Info}
	ErrMissingMajor          = &ModelError{"Please enter your major.", flash.Info}
	ErrMissingGraduationYear = &ModelError{"Please enter your graduation year.", flash.Info}
	ErrMissingGender         = &ModelError{"Please enter your gender.", flash.Info}
	ErrMissingGithub         = &ModelError{"Please enter your GitHub username.", flash.Info}
	ErrMissingPhone          = &ModelError{"Please enter your phone number.", flash.Info}
	ErrMissingWhyBM          = &ModelError{"Please enter why you want to come to BoilerMake.", flash.Info}
	ErrMissingLocation       = &ModelError{"Please enter your Location.", flash.Info}
	ErrMissingOtherSchool    = &ModelError{"Please enter your school's name.", flash.Info}
	ErrMissingOtherMajor     = &ModelError{"Please enter your major.", flash.Info}
	ErrMissingTACAgree       = &ModelError{"Please agree to the terms and conditions.", flash.Info}
	ErrMissingFirstName      = &ModelError{"Please enter your first name.", flash.Info}
	ErrMissingLastName       = &ModelError{"Please enter your last name.", flash.Info}

	// Validation errors when form paring
	ErrMissingResume  = &ModelError{"Please upload a resume.", flash.Info}
	ErrResumeTooLarge = &ModelError{"Resume upload is too large.", flash.Info}

	// Raffle errors
	ErrRaffleEmpty = &ModelError{"Please enter a raffle.", flash.Info}
	ErrInvalidRaffle = &ModelError{"That raffle code doesn't exist.", flash.Info}
	ErrTime = &ModelError{"This raffle has expired.", flash.Info}
	ErrRaffleClaimed = &ModelError{"You have already claimed this raffle.", flash.Info}
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
	ID              int
	Decision        int
	EmailedDecision bool
	UserID          int
	RSVP            bool
	AcceptedAt      time.Time
	CheckedInAt     time.Time

	School               string
	OtherSchool          string
	Major                string
	OtherMajor           string
	GraduationYear       string
	FirstName            string
	LastName             string
	ResumeFile           string
	Resume               *multipart.FileHeader // Stored in S3, not db
	Phone                string
	Gender               string
	Github               string
	Location             string
	IsFirstHackathon     bool
	WhyBM                string
	ProjIdea             string
	Is18OrOlder          bool
	MLHCodeOfConduct     bool
	MLHContestAndPrivacy bool
	Points               int
}

// Validate checks if an Application has all the necessary fields. Validation
// of resume uploads happens in application_service.go.
func (a *Application) Validate() error {
	if a.FirstName == "" {
		return ErrMissingFirstName
	} else if a.LastName == "" {
		return ErrMissingLastName
	} else if a.Gender == "" {
		return ErrMissingGender
	} else if a.Phone == "" {
		return ErrMissingPhone
	} else if a.Location == "" {
		return ErrMissingLocation
	} else if a.School == "" {
		return ErrMissingSchool
	} else if a.School == "Other" && a.OtherSchool == "" {
		return ErrMissingOtherSchool
	} else if a.Major == "" {
		return ErrMissingMajor
	} else if a.Major == "Other" && a.OtherMajor == "" {
		return ErrMissingOtherMajor
	} else if a.GraduationYear == "" {
		return ErrMissingGraduationYear
	} else if a.Github == "" {
		return ErrMissingGithub
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
	a.FirstName = r.FormValue("first-name")
	a.LastName = r.FormValue("last-name")
	a.Phone = r.FormValue("phone")
	a.Location = r.FormValue("location")
	a.School = r.FormValue("school")
	if a.School == "Other" { // only set OtherSchool if user selected "Other"
		a.OtherSchool = r.FormValue("other-school")
	}
	a.Major = r.FormValue("major")
	if a.Major == "Other" { // only set OtherMajor if user selected "Other"
		a.OtherMajor = r.FormValue("other-major")
	}
	a.GraduationYear = r.FormValue("graduation-year")
	a.Github = r.FormValue("github-username")
	a.Gender = r.FormValue("gender")
	a.WhyBM = r.FormValue("why-bm")
	a.ProjIdea = r.FormValue("proj-idea")

	// If no file was uploaded then set ResumeFile to empty string and let
	// application_service decide what to do.  If there's already a ResumeFile
	// in the db, then they've already uploaded a resume but just haven't updated
	// it with this post request.
	_, header, err := r.FormFile("resume")
	if err != nil {
		// Check if this error happened because request was too large
		// Kinda janky but it works
		if err.Error() == "multipart: NextPart: http: request body too large" {
			return ErrResumeTooLarge
		}

		// Otherwise only return an error if it's not because the file was missing.  Again,
		// we handle the missing resume case in the database.
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
	GetApplicationCount() int
	AddPointsToUser(uid int, points int) error
}
