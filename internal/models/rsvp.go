package models

import (
	"net/http"

	"github.com/BoilerMake/bm-app/pkg/flash"
)

var (
	ErrMissingShirtSize = &ModelError{"Please enter your shirt size.", flash.Info}
)

type RSVP struct {
	ID     int
	UserID int

	WillAttend     bool
	Accommodations string
	ShirtSize      string
	Allergies      string
}

// Validate checks if an RSVP has all the necessary fields.
func (rsvp *RSVP) Validate() error {
	if rsvp.ShirtSize == "" {
		return ErrMissingShirtSize
	}

	return nil
}

// FromFormData converts an RSVP from a requests's FormData to a
// models.RSVP struct.
func (rsvp *RSVP) FromFormData(r *http.Request) error {
	rsvp.WillAttend = r.FormValue("will-attend") == "on"
	rsvp.Accommodations = r.FormValue("accommodations")
	rsvp.ShirtSize = r.FormValue("shirt-size")
	rsvp.Allergies = r.FormValue("allergies")

	return nil
}

// An RSVPService defines an interface between the RSVP struct and its representation
// in our database.
type RSVPService interface {
	CreateOrUpdate(r *RSVP) error
	GetByUserID(uid int) (*RSVP, error)
}
