package models

import (
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// PayEvent model
type PayEvent struct {
	ID               IDType          `json:"id"`
	CreatedAt        time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time       `json:"updatedAt" db:"updated_at"`
	PlayedBy         *Person         `json:"playedby"`
	PlayedByID       *IDType         `json:"playedby_id" db:"playedby_id"`
	Beneficiary      *Person         `json:"beneficiary"`
	BeneficiaryID    *IDType         `json:"beneficiary_id" db:"beneficiary_id"`
	MusicRecording   *MusicRecording `json:"musicrecording"`
	MusicRecordingID *IDType         `json:"musicrecording_id" db:"musicrecording_id"`
	Amount           float64         `json:"amount"`
	Link             *string         `json:"link"`
}

func (pe *PayEvent) validate() error {
	if pe.PlayedByID == nil || *pe.PlayedByID == 0 {
		return errors.Errorf("PayEvent.PlayedByID is empty")
	} else if pe.MusicRecordingID == nil || *pe.MusicRecordingID == 0 {
		return errors.Errorf("PayEvent.MusicRecordingID is empty")
	}
	return nil
}

func (pe *PayEvent) Create(db *sqlx.DB) (id IDType, err error) {
	// err = pe.validate()
	// if err != nil {
	// 	return 0, err
	// }

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

	var payeventID IDType

	err = tx.QueryRow(
		`INSERT INTO payevent(
            created_at, updated_at, amount, link, playedby_id, beneficiary_id, musicrecording_id
        ) VALUES($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`,
		time.Now(), time.Now(), pe.Amount, pe.Link, pe.PlayedByID, pe.BeneficiaryID, pe.MusicRecordingID,
	).Scan(&payeventID)
	if err != nil {
		return 0, errors.Wrap(err, "error creating PayEvent")
	}
	return payeventID, nil
}

func (*PayEvent) Get(db *sqlx.DB, id IDType) (*PayEvent, error) {
	query := &SelectQuery{}
	query.SetQuery(`SELECT * FROM payevent WHERE payevent.id=$1`, id)

	var payevent PayEvent
	err := db.Get(&payevent, query.Query(), query.Args()...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting PayEvent")
	}

	if payevent.PlayedByID != nil && *payevent.PlayedByID > 0 {
		payevent.PlayedBy, err = (&Person{}).Get(db, *payevent.PlayedByID, "")
		if err != nil {
			lg.Errorf("error getting PayEvent.PlayedBy: %v", err)
		}
	}

	if payevent.BeneficiaryID != nil && *payevent.BeneficiaryID > 0 {
		payevent.Beneficiary, err = (&Person{}).Get(db, *payevent.BeneficiaryID, "")
		if err != nil {
			lg.Errorf("error getting PayEvent.Beneficiary: %v", err)
		}
	}

	if payevent.MusicRecordingID != nil && *payevent.MusicRecordingID > 0 {
		payevent.MusicRecording, err = (&MusicRecording{}).Get(db, *payevent.MusicRecordingID)
		if err != nil {
			lg.Errorf("error getting PayEvent.MusicRecording: %v", err)
		}
	}

	return &payevent, nil
}

func (*PayEvent) GetUnpaidForUser(db *sqlx.DB, userID *IDType, query *SelectQuery) ([]*PayEvent, error) {
	if query.Limit == 0 {
		query.Limit = 20
	}

	if userID == nil {
		query.SetQuery(`SELECT * FROM payevent WHERE beneficiary_id IS NULL AND link IS NOT NULL`)
	} else {
		query.SetQuery(`SELECT * FROM payevent WHERE beneficiary_id = $1 AND link IS NOT NULL`, userID)
	}

	payevents := []*PayEvent{}
	err := db.Select(&payevents, query.Query(), query.Args()...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting unpaid PayEvents")
	}
	return payevents, nil
}

func (*PayEvent) GetPaidForUser(db *sqlx.DB, userID IDType, query *SelectQuery) ([]*PayEvent, error) {
	if query.Limit == 0 {
		query.Limit = 20
	}

	query.SetQuery(`SELECT * FROM payevent WHERE beneficiary_id = $1 AND link IS NULL AND amount = 1`, userID)

	payevents := []*PayEvent{}
	err := db.Select(&payevents, query.Query(), query.Args()...)
	if err != nil {
		return nil, errors.Wrap(err, "error getting unpaid PayEvents")
	}
	return payevents, nil
}

func (*PayEvent) MarkAsPaid(db *sqlx.DB, ids []IDType, userID IDType) error {
	_, err := db.Exec(
		`UPDATE payevent
            SET link = NULL
            WHERE id = ANY($1) AND beneficiary_id = $2`,
		pq.Array(ids), userID,
	)
	if err != nil {
		return errors.Wrap(err, "error marking PayEvents as paid")
	}
	return nil
}

func (*PayEvent) MarkPrepaidAsClaimed(db *sqlx.DB, id IDType, userID IDType) error {
	_, err := db.Exec(
		`UPDATE payevent
            SET link = NULL, beneficiary_id = $1
            WHERE id = $2 AND beneficiary_id IS NULL`,
		userID, id,
	)
	if err != nil {
		return errors.Wrap(err, "error marking prepaid PayEvents as claimed")
	}
	return nil
}
