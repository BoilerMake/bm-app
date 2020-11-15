package models

import (
	"net/http"
	"time"

	"github.com/BoilerMake/bm-app/pkg/flash"
)

var (
	ErrMissingShirtSize = &ModelError{"Please enter your shirt size.", flash.Info}
	ErrExpired          = &ModelError{"Your RSVP has expired.  Contact team@boilermake.org if you think this was in error.", flash.Error}
)

const RSVPExpiryTime = 3 * 24 * time.Hour
const RSVPExpiryDate = 1608353999000 // this is milliseconds since epoch for the date 12/18/2020 23:59:59 EST

type RSVP struct {
	ID     int
	UserID int

	WillAttend bool
	OnCampus   bool
	ShirtSize  string
	StreetAddr string
	City       string
	State		string
	Country    string
	ZipCode    string
}

// Validate checks if an RSVP has all the necessary fields.
func (rsvp *RSVP) Validate() error {
	if rsvp.WillAttend && rsvp.ShirtSize == "" {
		return ErrMissingShirtSize
	}

	return nil
}

// FromFormData converts an RSVP from a requests's FormData to a
// models.RSVP struct.
func (rsvp *RSVP) FromFormData(r *http.Request) error {
	rsvp.WillAttend = r.FormValue("will-attend") == "on"
	rsvp.OnCampus = r.FormValue("on-campus") == "on"
	rsvp.ShirtSize = r.FormValue("shirt-size")
	rsvp.StreetAddr = r.FormValue("street-address")
	rsvp.City = r.FormValue("city")
	rsvp.State = r.FormValue("state")
	rsvp.Country = r.FormValue("country")
	rsvp.ZipCode = r.FormValue("zipcode")

	return nil
}

// An RSVPService defines an interface between the RSVP struct and its representation
// in our database.
type RSVPService interface {
	CreateOrUpdate(r *RSVP) error
	GetByUserID(uid int) (*RSVP, error)
}
