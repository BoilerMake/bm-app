package postgres

import (
	"database/sql"
	"strconv"

	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/lib/pq"
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

// toModel converts a database specific dbRaffle to the more generic
// Raffle struct
func (s *dbRaffle) toModel() *models.Raffle {
	return &models.Raffle{
		Code:      s.Code.String,
		StartTime: strconv.FormatInt(s.Start.Int64, 10),
		EndTime:   strconv.FormatInt(s.End.Int64, 10),
		Points:    strconv.FormatInt(s.Points.Int64, 10),
	}
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

	// Check for duplicate. Do this through checking postgres specific error // TODO: CHECK FOR EXEC CREATING RAFFLES
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "unique_violation":
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return models.ErrDuplicateRaffle
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return pgerr
		}
	}
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()
	return err
}

// Get raffle by id
func (s *RaffleService) GetById(id string) (*models.Raffle, error) {
	var dbr dbRaffle

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(`SELECT 
					raffle_id,
					start_time,
					end_time,
					points
				FROM raffles
				WHERE raffle_id = $1`, id).Scan(
		&dbr.Code,
		&dbr.Start,
		&dbr.End,
		&dbr.Points,
	)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	err = tx.Commit()
	return dbr.toModel(), err
}

// claim a raffle by inserting to raffle_hacker relation
func (s *RaffleService) ClaimRaffle(userId int, raffleId string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO raffle_hacker (user_id, raffle_id) VALUES ($1, $2);`,
		userId,
		raffleId,
	)

	// Check if user already claimed raffle. Do this through checking postgres specific error
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "unique_violation":
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return models.ErrRaffleClaimed
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return pgerr
		}
	}

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()
	return err
}
