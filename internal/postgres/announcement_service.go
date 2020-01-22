package postgres

import (
	"database/sql"

	"github.com/BoilerMake/bm-app/internal/models"
)

// AnnouncementService is a PostgreSQL implementation of models.announcement
type AnnouncementService struct {
	DB *sql.DB
}

// A dbAnnouncement is like a models.Announcement
// type dbAnnouncement struct {
// }

// Create creates an announcement and stores it in the DB
func (s *AnnouncementService) Create(message string) (err error) {
	// Get the DB
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// Create announcement with given message
	_, err = tx.Exec(`INSERT INTO announcements (message)
					  VALUES ($1);`,
		message,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()

	return err
}

// GetCurrent returns the most recent announcement (by timestamp)
func (s *AnnouncementService) GetCurrent() (a *models.Announcement, err error) {
	var dbA models.Announcement

	tx, err := s.DB.Begin()
	if err != nil {
		return a, err
	}

	err = tx.QueryRow(`SELECT
		id,
		message,
		created_at
	FROM announcements
	ORDER BY created_at DESC LIMIT 1`).Scan(&dbA.ID, &dbA.Message, &dbA.CreatedAt)

	if err != nil {
		// No announcements yet
		if err == sql.ErrNoRows {
			return nil, models.ErrNoAnnouncements
		}

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}

		return nil, err
	}

	err = tx.Commit()
	return &dbA, nil
}

// RemoveByID removes the announcement by id
// Note: removing by latest can cause racing call problems
// Should always remove by ID
func (s *AnnouncementService) DeleteByID(id int) (err error) {
	// Get the DB
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// Create announcement with given message
	_, err = tx.Exec(`DELETE FROM announcements WHERE id=$1;`,
		id,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()

	return err
}
