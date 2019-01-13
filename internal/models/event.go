// TODO I'm not really convinced this is the best way to structure our data
package models

/*
import "time"

type Announcement struct {
	ID        string    `json:"id"` // NOT NULL
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}

type Event struct {
	ID            string          `json:"id"` // NOT NULL
	Announcements []*Announcement `json:"announcements"`
	Applications  []*Application  `json:"applications"`
}

type EventService interface {
	Create(e *Event) error
	GetById(id string) (*Event, error)
	Update(e *Event) error

	CreateAnnouncement(a *Announcement) error
	UpdateAnnouncement(a *Announcement) error

	CreateApplication(a *Application) error
	UpdateApplication(a *Application) error
}
*/
