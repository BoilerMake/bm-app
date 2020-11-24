package postgres

import (
	"database/sql"
	"github.com/BoilerMake/bm-app/internal/models"
	"strconv"
	"strings"
)

// RaffleService is a PostgreSQL implementation of models.raffle
type RaffleService struct {
	DB *sql.DB
}

// A dbRaffle is like a models.Raffle, but it can read in null fields
// from the database without panicking
// need to convert between string and int with start,end,points
type dbRaffle struct {
	Code   sql.NullString
	Start  sql.NullInt64
	End    sql.NullInt64
	Points sql.NullInt64
}

// Create creates a raffle and stores it in the DB
func (s *RaffleService) Create(ra *models.Raffle) error {
	// Get the DB
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// Convert time and points from strings to integers
	startTimeInt, err := strconv.ParseInt(ra.StartTime, 10, 64)
	if err != nil {
		return models.ErrInvalidStartTime
	}
	endTimeInt, err := strconv.ParseInt(ra.EndTime, 10, 64)
	if err != nil {
		return models.ErrInvalidEndTime
	}
	pointsInt, err := strconv.Atoi(ra.Points)
	if err != nil {
		return models.ErrInvalidPoints
	}

	// Create raffle with entries
	_, err = tx.Exec(`INSERT INTO raffles (
			raffle_id,
			start_time,
			end_time,
			points
	)VALUES ($1, $2, $3, $4);`,
		ra.Code,
		startTimeInt,
		endTimeInt,
		pointsInt,
	)

	// Check for duplicate
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
