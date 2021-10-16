package models

import (
	"fmt"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/jmoiron/sqlx"
)

// Right model
type Right struct {
	ID               IDType    `json:"id"`
	CID              string    `json:"cid"`
	Type             string    `json:"@type"`
	Context          string    `json:"@context"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`
	ValidFrom        string    `json:"validFrom" db:"valid_from"`
	ValidThrough     string    `json:"validThrough" db:"valid_through"`
	PercentageShares float64   `json:"percentageShares" db:"percentage_shares"`
	Party            *Person   `json:"party"`
	PartyID          *IDType   `json:"-" db:"person_id"`
	// RightsOf         MusicRecording `json:"rightsOf"`
	// RightsOfID       *IDType        `json:"-" db:"musicrecording_id" // TODO - Do we need this?`
	// TODO - RightsOf could potentially refer to other types. Need to figure out how to handle that
	// with potentially an interface type and type assertion
	// TODO Add RightsType field
}

// GetParty gets the person object for a MusicRecording Right
// Note - This method replaced a call to Person.Get in order to also return musicgroup_member
// properties to the party object. Person.Get wouldn't work because a person isn't a musicgroup_member
// when they are first created unless they are the initial admin
func (obj *Right) GetParty(db *sqlx.DB, personID IDType) (*Person, error) {
	query := `
	SELECT person.*,
		musicgroup_members.description, musicgroup_members.percentage_shares,
		musicgroup_members.musicgroup_admin
	FROM person
	INNER JOIN musicgroup_members ON musicgroup_members.person_id = person.id
	WHERE person.id=$1
	`

	person := &Person{}
	err := db.Get(person, query, personID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = nil
		} else {
			return nil, fmt.Errorf("Error getting Right party (%v)", err)
		}
	}

	imageobject := &ImageObject{}
	imageobject, err = imageobject.Get(db, person.ImageID)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = nil
		} else {
			lg.Errorf("Error Scanning ImageObject in Person.Get (%v)", err)
		}
	}

	person.Image = imageobject

	return person, nil
}
