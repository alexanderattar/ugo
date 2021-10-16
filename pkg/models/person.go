package models

// Person model
import (
	"fmt"
	"regexp"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Person model
type Person struct {
	ID               IDType         `json:"id"`
	CID              string         `json:"cid"`
	CIDs             pq.StringArray `json:"cids"`
	Type             string         `json:"@type"`
	Context          string         `json:"@context"`
	CreatedAt        time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time      `json:"updatedAt" db:"updated_at"`
	EthereumAddress  string         `json:"ethereumAddress" db:"ethereum_address"`
	GivenName        *string        `json:"givenName" db:"given_name"`
	FamilyName       *string        `json:"familyName" db:"family_name"`
	Email            *string        `json:"email"`
	Image            *ImageObject   `json:"image"`
	ImageID          *int64         `json:"-" db:"image_id"`
	Description      *string        `json:"description"`
	PercentageShares *float64       `json:"percentageShares" db:"percentage_shares"`
	MusicGroupAdmin  *bool          `json:"musicgroupAdmin" db:"musicgroup_admin"`
	PaymentAddress   *string        `json:"paymentAddress" db:"payment_address"`
}

// All gets all of the given Person objects
func (obj *Person) All(db *sqlx.DB, ethereumAddress string, query *SelectQuery) ([]*Person, error) {
	if ethereumAddress != "" {
		query.SetQuery(`
            SELECT *
            FROM person
            WHERE person.ethereum_address=$1
        `, ethereumAddress)
	} else {
		query.SetQuery(`
            SELECT *
            FROM person
        `)
	}

	persons := []*Person{}
	err := db.Select(&persons, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting all Persons (%v)", err)
	}

	for _, person := range persons {
		imageobject := &ImageObject{}
		imageobject, err := imageobject.Get(db, person.ImageID)

		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				err = nil
			} else {
				return nil, fmt.Errorf("Error Scanning ImageObject in Person (%v)", err)
			}
		}

		person.Image = imageobject
	}

	return persons, err
}

// Get Person by ID or Eth Address
func (obj *Person) Get(db *sqlx.DB, personID IDType, ethereumAddress string) (*Person, error) {
	query := &SelectQuery{}
	if ethereumAddress == "" {
		query.SetQuery(`
            SELECT *
            FROM person
            WHERE person.id=$1
        `, personID)
	} else {
		query.SetQuery(`
            SELECT *
            FROM person
            WHERE person.ethereum_address=$1
        `, ethereumAddress)
		query.OrderBy = "person.created_at"
	}

	person := &Person{}
	err := db.Get(person, query.Query(), query.Args()...)
	if err != nil {
		// NOTE: Don't change error message unless you change the string
		// comparison in auth.go
		return nil, fmt.Errorf("%v", err)
	}

	imageobject := &ImageObject{}
	imageobject, err = imageobject.Get(db, person.ImageID)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = nil
		} else {
			lg.Errorf("Error Scanning ImageObject in Person.Get (%v)", err)
		}
	} else {
		person.Image = imageobject
	}
	return person, nil
}

// GetByCID Person by CID
func (obj *Person) GetByCID(db *sqlx.DB, args ...interface{}) (*Person, error) {
	query := `
	SELECT *
	FROM person
	WHERE $1 <@ cids;
	`

	person := &Person{}
	err := db.Get(person, query, pq.Array(args))

	if err != nil {
		return nil, fmt.Errorf("Error getting Person (%v)", err)
	}

	imageobject := &ImageObject{}
	imageobject, err = imageobject.Get(db, person.ImageID)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = nil
		} else {
			lg.Errorf("Error Scanning ImageObject in Person.Get (%v)", err)
		}
	} else {
		person.Image = imageobject
	}

	person.Image = imageobject

	return person, nil
}

// Create a Person
func (obj *Person) Create(db *sqlx.DB) (personID IDType, err error) {
	err = obj.validate()
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

	var CIDs []string
	// On initial POST from portal, there is no CID
	if obj.CID != "" {
		CIDs = []string{obj.CID}
	}

	var imageID interface{}

	if obj.Image != nil {
		// Create ImageObject
		err = tx.QueryRow(
			`INSERT INTO imageobject (
				cid, type, context, created_at, updated_at,
				content_url, encoding_format
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`,
			obj.Image.CID, obj.Image.Type, obj.Image.Context, time.Now(), time.Now(),
			obj.Image.ContentURL, obj.Image.EncodingFormat,
		).Scan(&imageID)

		if err != nil {
			return 0, fmt.Errorf("Error creating person image (%v)", err)
		}
	}

	// Create Person
	err = tx.QueryRow(
		`INSERT INTO person(
			cid, cids, type, context, created_at, updated_at,
			ethereum_address, given_name, family_name, email, image_id, payment_address
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`,
		obj.CID, pq.Array(CIDs), obj.Type, obj.Context, time.Now(), time.Now(),
		obj.EthereumAddress, obj.GivenName, obj.FamilyName, obj.Email, imageID,
		obj.PaymentAddress,
	).Scan(&personID)

	if err != nil {
		return 0, fmt.Errorf("Error creating Person (%v)", err)
	}
	return personID, nil
}

// Update a Person
func (obj *Person) Update(db *sqlx.DB, id IDType) (personID IDType, err error) {
	err = obj.validate()
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

	if obj.Image == nil {
		// Update person without ImageObject
		// We can remove this if an ImageObject is required by Person
		err = tx.QueryRow(
			`UPDATE person
			 SET cid=$1, cids=array_append(cids, $2), type=$3, context=$4, updated_at=$5,
			 ethereum_address=$6, given_name=$7, family_name=$8, email=$9, payment_address=$10
			 WHERE id=$11
			 RETURNING id`,
			obj.CID, obj.CID, obj.Type, obj.Context, time.Now(),
			obj.EthereumAddress, obj.GivenName, obj.FamilyName, obj.Email, obj.PaymentAddress,
			id,
		).Scan(&personID)

		if err != nil {
			return personID, fmt.Errorf("Error updating person (%v)", err)
		}
	} else if obj.Image.ID < 1 {
		// Create ImageObject
		_, err = tx.Exec(
			`INSERT INTO imageobject (
				cid, type, context, created_at, updated_at,
				content_url, encoding_format
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			obj.Image.CID, obj.Image.Type, obj.Image.Context, time.Now(), time.Now(),
			obj.Image.ContentURL, obj.Image.EncodingFormat,
		)

		if err != nil {
			return personID, fmt.Errorf("Error updating image in person (%v)", err)
		}

		// Update Person with new ImageObject
		err = tx.QueryRow(
			`UPDATE person
			 SET cid=$1,  cids=array_append(cids, $2), type=$3, context=$4, updated_at=$5,
			 ethereum_address=$6, given_name=$7, family_name=$8, email=$9,
			 image_id=currval('imageobject_id_seq'), payment_address=$10
			 WHERE id=$11
			 RETURNING id`,
			obj.CID, obj.CID, obj.Type, obj.Context, time.Now(),
			obj.EthereumAddress, obj.GivenName, obj.FamilyName, obj.Email, obj.PaymentAddress,
			id,
		).Scan(&personID)

		if err != nil {
			return personID, fmt.Errorf("Error updating person (%v)", err)
		}
	} else if obj.Image.ID > 0 {
		// Update Image
		_, err = tx.Exec(
			`UPDATE imageobject
			 SET cid=$1, type=$2, context=$3,
			 updated_at=$4, content_url=$5, encoding_format=$6
			 WHERE id=$7`,
			obj.Image.CID, obj.Type, obj.Context, time.Now(),
			obj.Image.ContentURL, obj.Image.EncodingFormat, obj.Image.ID,
		)

		if err != nil {
			return personID, fmt.Errorf("Error updating image in person (%v)", err)
		}

		// Update Person
		err = tx.QueryRow(
			`UPDATE person
			 SET cid=$1, cids=array_append(cids, $2), type=$3, context=$4, updated_at=$5,
			 ethereum_address=$6, given_name=$7, family_name=$8, email=$9, payment_address=$10
			 WHERE id=$11
			 RETURNING id`,
			obj.CID, obj.CID, obj.Type, obj.Context, time.Now(),
			obj.EthereumAddress, obj.GivenName, obj.FamilyName, obj.Email, obj.PaymentAddress,
			id,
		).Scan(&personID)

		if err != nil {
			return 0, fmt.Errorf("Error updating person (%v)", err)
		}
	}
	return personID, nil
}

// Delete a Person
func (obj *Person) Delete(db *sqlx.DB, personID IDType) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = db.Exec(
		`DELETE FROM person
         WHERE person.id = $1`,
		personID,
	)

	if err != nil {
		return fmt.Errorf("Error deleting person (%v)", err)
	}
	return nil
}

func (obj Person) validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.EthereumAddress, validation.Required, validation.Match(regexp.MustCompile("^0x[0-9a-fA-F]{40}$"))),
		validation.Field(&obj.Email, is.Email),
	)
}
