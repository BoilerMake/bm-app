package models

import (
	"github.com/BoilerMake/bm-app/pkg/flash"
	"net/http"
)

// Raffle errors
var (
	ErrRaffleCodeEmpty  = &ModelError{"Raffle code is empty", flash.Error}
	ErrStartTimeEmpty   = &ModelError{"Start Time is missing", flash.Error}
	ErrInvalidStartTime = &ModelError{"Incorrect format for start time", flash.Error}
	ErrEndTimeEmpty     = &ModelError{"End Time is missing", flash.Error}
	ErrInvalidEndTime   = &ModelError{"Incorrect format for end time", flash.Error}
	ErrPointsEmpty      = &ModelError{"Points is missing", flash.Error}
	ErrInvalidPoints    = &ModelError{"Incorrect format for points", flash.Error}
	ErrDuplicateRaffle  = &ModelError{"This raffle code already exists", flash.Error}

	ErrInvalidPointsToAdd = &ModelError{"Incorrect format for points to add", flash.Error}
)

// A Raffle is a raffle stored in the raffles table
type Raffle struct {
	Code string `json:"string"`

	// parse as string and convert to int in db
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Points    string `json:"points"`
}

// Validate checks if a raffle has all necessary fields
func (ra *Raffle) Validate() error {
	if ra.Code == "" {
		return ErrRaffleCodeEmpty
	} else if ra.StartTime == "" {
		return ErrStartTimeEmpty
	} else if ra.EndTime == "" {
		return ErrEndTimeEmpty
	} else if ra.Points == "" {
		return ErrPointsEmpty
	}

	return nil
}

// FromFormData builds raffle struct
func (ra *Raffle) FromFormData(r *http.Request) error {
	ra.Code = r.FormValue("code")
	ra.StartTime = r.FormValue("starttime")
	ra.EndTime = r.FormValue("endtime")
	ra.Points = r.FormValue("points")

	return nil
}

// A RaffleService defines interface between the
// rsvp model and database representation of the raffles table
type RaffleService interface {
	Create(ra *Raffle) error
	GetById(id string) (*Raffle, error)
	ClaimRaffle(userId int, raffleId string) error
}
