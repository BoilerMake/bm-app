package models

import (
	"time"

	"github.com/BoilerMake/bm-app/pkg/flash"
)

// Authentication errors
var (
	ErrAnnouncementMessageEmpty = &ModelError{"Announcement message is empty.", flash.Error}
	ErrAnnouncementIDEmpty      = &ModelError{"Announcement ID is empty.", flash.Error}
	ErrNoAnnouncements          = &ModelError{"No announcements in database.", flash.Info}
	ErrAnnouncementNotFound     = &ModelError{"Announcement not found.", flash.Info}
)

// An Announcement is an announcement stored in the database.
type Announcement struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

// An AnnouncementService defines an interface between the
// announcement model and the db representation.
type AnnouncementService interface {
	Create(message string) error
	GetByID(id int) (*Announcement, error)
	GetCurrent() (*Announcement, error)
	DeleteByID(id int) error
}
