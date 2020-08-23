package models

import (
	"github.com/BoilerMake/bm-app/pkg/flash"
	"net/http"
	"strconv"
)

// Raffle errors
var (
	ErrRaffleCodeEmpty = &ModelError{"Raffle code is empty", flash.Error}
	ErrInvalidStartTime = &ModelError{"Invalid Start Time Format", flash.Error}
	// ErrBigStartTime = &ModelError{"Start time to large", flash.Error}
	ErrInvalidEndtime = &ModelError{"Invalid End Time Format", flash.Error}
	// ErrBigEndTime = &ModelError{"End time to large", flash.Error}
	ErrInvalidPoints = &ModelError{"Invalid Points Format", flash.Error}
	// ErrPointsToBig = &ModelError{"Points out of bounds(int32)", flash.Error}
)

// A Raffle is an raffle stored in the database.
type Raffle struct {
	Code		string	`json:"string"`
	StartTime	int64	`json:"startTime"`
	EndTime		int64	`json:"endTime"`
	Points		int		`json:"points"`
}

// Validate checks if a raffle has all necessary fields.
func (ra *Raffle) Validate() error { // take in input from user entered string
	return nil
}

// FromFormData converts exec input and validates to make sure int is in place
func (ra *Raffle) FromFormData(r *http.Request) error {
	ra.Code = r.FormValue("code")
	if ra.Code == "" {
		return ErrRaffleCodeEmpty
	}
	start, err := strconv.ParseInt(r.FormValue("starttime"), 10, 64)
	if err != nil { // empty string and invalid string has same error message
		return ErrInvalidStartTime
	}
	ra.StartTime = start
	end, err := strconv.ParseInt(r.FormValue("endtime"), 10, 64)
	if err != nil {
		return ErrInvalidEndtime
	}
	ra.EndTime = end
	points, err := strconv.ParseInt(r.FormValue("points"), 10, 32)
	if err != nil {
		return ErrInvalidPoints
	}
	ra.Points = int(points) // strconv.ParseInt returns an int64 regardless
	return nil
}

// An AnnouncementService defines an interface between the
// announcement model and the db representation.
type RaffleService interface {
	Create(ra *Raffle) error
	GetByCode(code string) (*Raffle, error)
}
