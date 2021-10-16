package models

import (
	"fmt"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/jmoiron/sqlx"
)

// Report model
type Report struct {
	ID             IDType        `json:"id"`
	CreatedAt      time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time     `json:"updatedAt" db:"updated_at"`
	State          string        `json:"state"`
	Response       *string       `json:"response"`
	Reason         string        `json:"reason"`
	Message        *string       `json:"message"`
	Email          *string       `json:"email"`
	MusicRelease   *MusicRelease `json:"musicrelease"`
	ReporterID     IDType        `json:"reporter_id" db:"reporter_id"`
	MusicReleaseID IDType        `json:"musicrelease_id" db:"musicrelease_id"`
}

// Create a Report
func (obj *Report) Create(db *sqlx.DB) (id IDType, err error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	obj.State = "unreviewed"

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Create Report
	var reportID IDType
	err = tx.QueryRow(`
        INSERT INTO report(
        created_at, updated_at, state, reason, message,
        email, musicrelease_id, reporter_id)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id
    `,
		time.Now(), time.Now(), obj.State, obj.Reason, obj.Message, obj.Email,
		obj.MusicReleaseID, obj.ReporterID,
	).Scan(&reportID)

	if err != nil {
		return 0, fmt.Errorf("Error inserting into report (%v)", err)
	}

	musicrelease, err := (&MusicRelease{}).GetByID(db, obj.MusicReleaseID)
	if err != nil {
		lg.Errorf("Error setting music release on report (%v)", err)
	}

	obj.MusicRelease = musicrelease

	return reportID, nil
}

// All gets all of the given Report objects
func (obj *Report) All(db *sqlx.DB, musicreleaseID *IDType, query *SelectQuery) ([]*Report, error) {
	if musicreleaseID != nil {
		query.SetQuery(`
            SELECT *
            FROM report
            WHERE report.musicrelease_id=$1
        `, *musicreleaseID)
	} else {
		query.SetQuery(`
            SELECT *
            FROM report
        `)
	}
	reports := []*Report{}
	err := db.Select(&reports, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting all Reports (%v)", err)
	}

	for _, report := range reports {
		musicrelease, err := (&MusicRelease{}).GetByID(db, report.MusicReleaseID)
		if err != nil {
			lg.Errorf("Error setting music release on report (%v)", err)
		}

		report.MusicRelease = musicrelease
	}

	return reports, nil
}

// Get Report by ID
func (obj *Report) Get(db *sqlx.DB, reportID IDType) (*Report, error) {
	query := &SelectQuery{}
	query.SetQuery(`
        SELECT *
        FROM report
        WHERE report.id=$1
    `, reportID)

	report := &Report{}
	err := db.Get(report, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting Report (%v)", err)
	}

	musicrelease, err := (&MusicRelease{}).GetByID(db, report.MusicReleaseID)
	if err != nil {
		lg.Errorf("Error setting music release on report (%v)", err)
	}

	report.MusicRelease = musicrelease

	return report, nil
}

// Update a Report
func (obj *Report) Update(db *sqlx.DB, reportID IDType) (r *Report, err error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Update Report
	err = tx.QueryRow(
		`UPDATE report
         SET created_at=$1, updated_at=$2, state=$3, reason=$4, message=$5,
         response=$6, email=$7, musicrelease_id=$8, reporter_id=$9
         WHERE id=$10
         RETURNING id`,
		time.Now(), time.Now(), obj.State, obj.Reason, obj.Message, obj.Response,
		obj.Email, obj.MusicReleaseID, obj.ReporterID,
		reportID,
	).Scan(&reportID)

	if err != nil {
		return nil, fmt.Errorf("Error updating report (%v)", err)
	}

	musicrelease, err := (&MusicRelease{}).GetByID(db, obj.MusicReleaseID)
	if err != nil {
		lg.Errorf("Error setting music release on report (%v)", err)
	}

	obj.MusicRelease = musicrelease

	return obj, nil
}

// Resolve a Report
func (obj *Report) Resolve(db *sqlx.DB, reportID IDType) (*Report, error) {
	obj.State = "resolved"

	return obj.Update(db, reportID)
}

// Deactivate a Report
func (obj *Report) Deactivate(db *sqlx.DB, reportID IDType) (*Report, error) {
	obj.State = "deactivated"

	musicrelease, err := (&MusicRelease{}).GetByID(db, obj.MusicReleaseID)
	if err != nil {
		return nil, err
	}

	musicrelease.Active = false
	musicrelease, err = musicrelease.Update(db, obj.MusicReleaseID)
	if err != nil {
		return nil, err
	}

	return obj.Update(db, reportID)
}
