package models

import (
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// PlayEvent model
type PlayEvent struct {
	ID               IDType          `json:"id"`
	CreatedAt        time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time       `json:"updatedAt" db:"updated_at"`
	PlayedBy         *Person         `json:"playedby"`
	PlayedByID       *IDType         `json:"playedby_id" db:"playedby_id"`
	MusicRecording   *MusicRecording `json:"musicrecording"`
	MusicRecordingID IDType          `json:"musicrecording_id" db:"musicrecording_id"`
}

func (pe *PlayEvent) validate() error {
	if pe.PlayedByID == nil || *pe.PlayedByID == 0 {
		return errors.Errorf("PlayEvent.PlayedByID is empty")
	} else if pe.MusicRecordingID == 0 {
		return errors.Errorf("PlayEvent.MusicRecordingID is empty")
	}
	return nil
}

func (pe *PlayEvent) Create(db *sqlx.DB) (IDType, error) {
	err := pe.validate()
	if err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var playeventID IDType

	err = tx.QueryRow(
		`INSERT INTO playevent(
            created_at, updated_at, playedby_id, musicrecording_id
        ) VALUES($1, $2, $3, $4)
        RETURNING id`,
		time.Now(), time.Now(), pe.PlayedByID, pe.MusicRecordingID,
	).Scan(&playeventID)
	if err != nil {
		return 0, errors.Wrap(err, "error creating PlayEvent")
	}
	return playeventID, nil
}

func (pe *PlayEvent) Get(db *sqlx.DB, id IDType) (*PlayEvent, error) {
	query := &SelectQuery{}
	query.SetQuery(`SELECT * FROM playevent WHERE playevent.id=$1`, id)

	var playevent PlayEvent
	err := db.Get(&playevent, query.Query(), query.Args()...)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting PlayEvent")
	}

	if playevent.PlayedByID != nil && *playevent.PlayedByID > 0 {
		playevent.PlayedBy, err = (&Person{}).Get(db, *playevent.PlayedByID, "")
		if err != nil {
			lg.Errorf("error getting PayEvent.PlayedBy: %v", err)
		}
	}

	if playevent.MusicRecordingID > 0 {
		playevent.MusicRecording, err = (&MusicRecording{}).Get(db, playevent.MusicRecordingID)
		if err != nil {
			lg.Errorf("error getting PayEvent.MusicRecording: %v", err)
		}
	}

	return &playevent, nil
}
