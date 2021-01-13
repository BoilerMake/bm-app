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

// GetByID returns an announcement with the given ID.
func (s *AnnouncementService) GetByID(id int) (a *models.Announcement, err error) {
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
	WHERE id = $1`, id).Scan(&dbA.ID, &dbA.Message, &dbA.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// Announcement with given ID doesn't exist
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

func (s *AnnouncementService) GetAllAnnouncements() (a []*models.Announcement, err error) {
	var dbAs []*models.Announcement
	tx, err := s.DB.Begin()
	if err != nil {
		return a, err
	}

	rows, err := tx.Query(`SELECT id, message, created_at FROM announcements ORDER BY created_at DESC`)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}

		return nil, err
	}

	counter := 0 // check if loop is ran
	defer rows.Close()
	for rows.Next() {
		var dbA models.Announcement
		err = rows.Scan(&dbA.ID, &dbA.Message, &dbA.CreatedAt)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}

		dbAs = append(dbAs, &dbA)
		counter += 1 // indicate that at least one row was found
	}

	err = rows.Err() // check if there was an error within loop
	if err != nil {
		return nil, err
	}
	if counter == 0 { // if the loop was never ran that means there was no rows returned
		return nil, models.ErrNoAnnouncements
	}

	err = tx.Commit()

	return dbAs, nil
}
