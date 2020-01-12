package postgres

import (
	"database/sql"

	"github.com/BoilerMake/bm-app/internal/models"
)

// RSVPService is a PostgreSQL implementation of models.RSVPService
type RSVPService struct {
	DB *sql.DB
}

// A dbRSVP is like a models.RSVP, but it can read in null fields
// from the database without panicking.
type dbRSVP struct {
	ID     sql.NullInt64
	UserID sql.NullInt64

	WillAttend     sql.NullBool
	Accommodations sql.NullString
	ShirtSize      sql.NullString
	Allergies      sql.NullString
}

// toModel converts a database specific dbRSVP to the more generic
// RSVP struct.
func (r *dbRSVP) toModel() *models.RSVP {
	return &models.RSVP{
		ID:     int(r.ID.Int64),
		UserID: int(r.UserID.Int64),

		WillAttend:     r.WillAttend.Bool,
		Accommodations: r.Accommodations.String,
		ShirtSize:      r.ShirtSize.String,
		Allergies:      r.Allergies.String,
	}
}

// CreateOrUpdate tries to make a new RSVP, unless one already exists then
// it updates the existing one.
func (s *RSVPService) CreateOrUpdate(r *models.RSVP) (err error) {
	err = r.Validate()
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// See if they already have an RSVP
	_, err = s.GetByUserID(r.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			// RSVP hasn't been made yet
			_, err = tx.Exec(`INSERT INTO rsvps (
			user_id,
			will_attend,
			accommodations,
			shirt_size,
			allergies
		) VALUES ($1, $2, $3, $4, $5);`,
				r.UserID,
				r.WillAttend,
				r.Accommodations,
				r.ShirtSize,
				r.Allergies,
			)

			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return rollbackErr
				}
				return err
			}
		} else {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	} else {
		// RSVP already exists, so update it
		_, err = tx.Exec(`UPDATE rsvps
		SET
			will_attend = $1,
			accommodations = $2,
			shirt_size = $3,
			allergies = $4
		WHERE user_id = $5`,
			r.WillAttend,
			r.Accommodations,
			r.ShirtSize,
			r.Allergies,
			r.UserID,
		)

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	err = tx.Commit()
	return err
}

// GetByUserID returns a single RSVP with the given user id.
func (s *RSVPService) GetByUserID(uid int) (*models.RSVP, error) {
	var dbr dbRSVP

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(`SELECT
			id,
			user_id,
			will_attend,
			accommodations,
			shirt_size,
			allergies
		FROM rsvps
		WHERE user_id = $1`, uid).Scan(
		&dbr.ID,
		&dbr.UserID,
		&dbr.WillAttend,
		&dbr.Accommodations,
		&dbr.ShirtSize,
		&dbr.Allergies,
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
