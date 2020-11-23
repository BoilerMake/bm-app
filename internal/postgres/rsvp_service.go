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

	WillAttend sql.NullBool
	OnCampus   sql.NullBool
	ShirtSize  sql.NullString
	StreetAddr sql.NullString
	City       sql.NullString
	State		sql.NullString
	Country    sql.NullString
	ZipCode    sql.NullString
}

// toModel converts a database specific dbRSVP to the more generic
// RSVP struct.
func (r *dbRSVP) toModel() *models.RSVP {
	return &models.RSVP{
		ID:     int(r.ID.Int64),
		UserID: int(r.UserID.Int64),

		WillAttend: r.WillAttend.Bool,
		OnCampus:   r.OnCampus.Bool,
		ShirtSize:  r.ShirtSize.String,
		StreetAddr: r.StreetAddr.String,
		City:       r.City.String,
		State:		r.State.String,
		Country:    r.Country.String,
		ZipCode:    r.ZipCode.String,
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
			shirt_size,
			oncampus,
			street_address,
			city,
			state,
			country,
			zip_code
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
				r.UserID,
				r.WillAttend,
				r.ShirtSize,
				r.OnCampus,
				r.StreetAddr,
				r.City,
				r.State,
				r.Country,
				r.ZipCode,
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
			shirt_size = $2,
			oncampus = $3,
			street_address = $4,
			city = $5,
			state = $6,
			country = $7,
			zip_code = $8
		WHERE user_id = $9`,
			r.WillAttend,
			r.ShirtSize,
			r.OnCampus,
			r.StreetAddr,
			r.City,
			r.State,
			r.Country,
			r.ZipCode,
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
			shirt_size,
			oncampus,
			street_address,
			city,
			state,
			country,
			zip_code
		FROM rsvps
		WHERE user_id = $1`, uid).Scan(
		&dbr.ID,
		&dbr.UserID,
		&dbr.WillAttend,
		&dbr.ShirtSize,
		&dbr.OnCampus,
		&dbr.StreetAddr,
		&dbr.City,
		&dbr.State,
		&dbr.Country,
		&dbr.ZipCode,
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
