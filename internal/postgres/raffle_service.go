package postgres

import (
	"database/sql"
	"github.com/BoilerMake/bm-app/internal/models"
	"strings"
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
		if strings.Contains(err.Error(), "duplicate key") {
			return models.ErrDuplicateRaffle
		}
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()
	return err
}

// GetByCode returns a raffle with the given code.
func (s *RaffleService) GetByCode(code string) (ra *models.Raffle, err error) {
	var dba models.Raffle

	tx, err := s.DB.Begin()
	if err != nil {
		return ra, err
	}

	err = tx.QueryRow(`SELECT
		code,
		start_time,
		end_time,
		points
	FROM raffles
	WHERE code=$1`, code).Scan(&dba.Code, &dba.StartTime, &dba.EndTime, &dba.Points)

	if err != nil {
		if err == sql.ErrNoRows {
			// Raffle with this code doesn't exist
			return nil, models.ErrInvalidRaffle
		}

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}

		return nil, err
	}

	err = tx.Commit()
	return &dba, nil
}

