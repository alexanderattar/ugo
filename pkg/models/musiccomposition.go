package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// MusicComposition model
type MusicComposition struct {
	ID        IDType    `json:"id"`
	CID       string    `json:"cid"`
	Type      string    `json:"@type"`
	Context   string    `json:"@context"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	Name      *string   `json:"name"`
	Composer  []*Person `json:"composer"`
	Iswc      *string   `json:"iswc"`
}

// GetComposer returns all the composers of a MusicComposition
func (obj *MusicComposition) GetComposer(db *sqlx.DB, args ...interface{}) ([]*Person, error) {
	query := `
	SELECT person.id, person.cid, person.type,
		person.context, person.created_at, person.updated_at, person.ethereum_address,
		person.given_name, person.family_name
	FROM person
	INNER JOIN musiccomposition_composer ON musiccomposition_composer.person_id = person.id
	WHERE musiccomposition_composer.musiccomposition_id=$1
	`

	persons := []*Person{}
	err := db.Select(&persons, query, args...)

	if err != nil {
		return nil, fmt.Errorf("Error getting MusicGroup members (%v)", err)
	}

	return persons, nil
}
