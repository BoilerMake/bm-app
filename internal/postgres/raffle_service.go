package postgres

import (
	"database/sql"
	"github.com/BoilerMake/bm-app/internal/models"
)

// AnnouncementService is a PostgreSQL implementation of models.announcement
type RaffleService struct {
	DB *sql.DB
}

// A dbRaffle is like a models.Raffle, but it can read in null fields
// from the database without panicking
type dbRaffle struct {
	Code	sql.NullString
	Start	sql.NullInt64
	End		sql.NullInt64
	Points	sql.NullInt64
}

// Create creates an announcement and stores it in the DB
func (s *RaffleService) Create(ra *models.Raffle) error {
	// Get the DB
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	// Create announcement with given message
	_, err = tx.Exec(`INSERT INTO raffles (
		code,
		start_time,
		end_time,
		points
	) VALUES ($1, $2, $3, $4);`,
			ra.Code,
			ra.StartTime,
			ra.EndTime,
			ra.Points,
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
func (s *RaffleService) GetByCode(code string) (a *models.Raffle, err error) {
	var dba models.Raffle
	return &dba, nil
}

