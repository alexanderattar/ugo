package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// SignedMessage model
type SignedMessage struct {
	ID              IDType    `json:"id"`
	Message         string    `json:"message"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
	EthereumAddress string    `json:"ethereumAddress" db:"ethereum_address"`
	Signature       string    `json:"signature" db:"signature"`
}

// All gets all of the given SignedMessage objects
func (obj *SignedMessage) All(db *sqlx.DB, query *SelectQuery) ([]*SignedMessage, error) {
	query.SetQuery(`
        SELECT *
        FROM signedmessage
    `)

	signedmessages := []*SignedMessage{}
	err := db.Select(&signedmessages, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting all SignedMessages (%v)", err)
	}

	return signedmessages, err
}

// Get SignedMessage by ID
func (obj *SignedMessage) Get(db *sqlx.DB, id IDType) (*SignedMessage, error) {
	query := &SelectQuery{}
	query.SetQuery(`
        SELECT *
        FROM signedmessage
        WHERE signedmessage.id=$1
    `, id)

	signedmessage := &SignedMessage{}
	err := db.Get(signedmessage, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting SignedMessage (%v)", err)
	}

	return signedmessage, err
}

// GetBySignedMessage gets by SignedMessage
func (obj *SignedMessage) GetBySignedMessage(db *sqlx.DB, signedMessage string) (*SignedMessage, error) {
	query := &SelectQuery{}
	query.SetQuery(`
        SELECT *
        FROM signedmessage
        WHERE signedmessage.signedmessage=$1
    `, signedMessage)

	signedmessage := &SignedMessage{}
	err := db.Get(signedmessage, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting SignedMessage (%v)", err)
	}

	return signedmessage, err
}

// Create a SignedMessage
func (obj *SignedMessage) Create(db *sqlx.DB) (id IDType, err error) {
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

	var signedmessageID IDType

	err = tx.QueryRow(
		`INSERT INTO signedmessage(
         message, created_at, updated_at,
         ethereum_address, signature
     ) VALUES($1, $2, $3, $4, $5)
     RETURNING id`,
		obj.Message, time.Now(), time.Now(),
		obj.EthereumAddress, obj.Signature,
	).Scan(&signedmessageID)

	if err != nil {
		return 0, fmt.Errorf("Error creating SignedMessage (%v)", err)
	}

	return signedmessageID, nil
}
