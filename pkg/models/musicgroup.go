package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// MusicGroup model
type MusicGroup struct {
	ID          IDType         `json:"id"`
	CID         string         `json:"cid"`
	CIDs        pq.StringArray `json:"cids"`
	Type        string         `json:"@type"`
	Context     string         `json:"@context"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at"`
	Name        *string        `json:"name"`
	Description *string        `json:"description"`
	Email       *string        `json:"-"`
	Members     []*Person      `json:"members"`
	Image       *ImageObject   `json:"image"`
	ImageID     *IDType        `json:"-" db:"image_id"`
}

// All gets all of the given MusicGroup objects
func (obj *MusicGroup) All(db *sqlx.DB, ethereumAddress string, personID *IDType, query *SelectQuery) ([]*MusicGroup, error) {
	if ethereumAddress != "" {
		query.SetQuery(`
        SELECT musicgroup.id, musicgroup.cid, musicgroup.cids, musicgroup.type,
          musicgroup.context, musicgroup.created_at, musicgroup.updated_at,
          musicgroup.name, musicgroup.description, musicgroup.email, musicgroup.image_id
        FROM musicgroup
        INNER JOIN musicgroup_members ON musicgroup_members.musicgroup_id = musicgroup.id
        INNER JOIN person ON musicgroup_members.person_id = person.id
        WHERE person.ethereum_address=$1
        `, ethereumAddress)

	} else if personID != nil {
		query.SetQuery(`
        SELECT musicgroup.id, musicgroup.cid, musicgroup.cids, musicgroup.type,
          musicgroup.context, musicgroup.created_at, musicgroup.updated_at,
          musicgroup.name, musicgroup.description, musicgroup.email, musicgroup.image_id
        FROM musicgroup
        INNER JOIN musicgroup_members ON musicgroup_members.musicgroup_id = musicgroup.id
        INNER JOIN person ON musicgroup_members.person_id = person.id
        WHERE person.id=$1
        `, *personID)

	} else {
		query.SetQuery(`
        SELECT *
        FROM musicgroup
        `)
	}

	musicgroups := []*MusicGroup{}
	err := db.Select(&musicgroups, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting all MusicGroups (%v)", err)
	}

	// TODO - Optimize query to grab all fields instead of round trips to the db for subobjects
	// Serialize the image
	for _, musicgroup := range musicgroups {
		image, err := musicgroup.GetImage(db, *musicgroup.ImageID)
		if err != nil {
			lg.Errorf("Error Scanning ImageObject (%v)", err)
		}

		musicgroup.Image = image
	}

	return musicgroups, nil
}

// AllByPersonID gets all of the given MusicGroup objects for person by personID
func (obj *MusicGroup) AllByPersonID(db *sqlx.DB, personID IDType, query *SelectQuery) ([]*MusicGroup, error) {
	query.SetQuery(`
		SELECT musicgroup.id, musicgroup.cid, musicgroup.cids, musicgroup.type,
			musicgroup.context, musicgroup.created_at, musicgroup.updated_at,
			musicgroup.name, musicgroup.description, musicgroup.email, musicgroup.image_id
		FROM musicgroup
		INNER JOIN musicgroup_members ON musicgroup_members.musicgroup_id = musicgroup.id
		INNER JOIN person ON musicgroup_members.person_id = person.id
		WHERE person.id=$1
		`, personID)

	musicgroups := []*MusicGroup{}
	err := db.Select(&musicgroups, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting all MusicGroups (%v)", err)
	}

	// TODO - Optimize query to grab all fields instead of round trips to the db for subobjects
	// Serialize the image
	for _, musicgroup := range musicgroups {
		image, err := musicgroup.GetImage(db, *musicgroup.ImageID)
		if err != nil {
			lg.Errorf("Error Scanning ImageObject (%v)", err)
		}

		musicgroup.Image = image
	}

	return musicgroups, nil
}

// Get MusicGroup by ID
func (obj *MusicGroup) Get(db *sqlx.DB, id IDType) (*MusicGroup, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT *
    FROM musicgroup
    WHERE musicgroup.id=$1
    `, id)

	musicgroup := &MusicGroup{}
	err := db.Get(musicgroup, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting MusicGroup (%v)", err)
	}

	members, err := musicgroup.GetMembers(db, musicgroup.ID)
	if err != nil {
		lg.Errorf("Error Scanning MembersObject (%v)", err)
	}

	image, err := musicgroup.GetImage(db, *musicgroup.ImageID)
	if err != nil {
		lg.Errorf("Error Scanning ImageObject (%v)", err)
	}

	musicgroup.Members = members
	musicgroup.Image = image

	return musicgroup, nil
}

// GetByCID MusicGroup by CID
func (obj *MusicGroup) GetByCID(db *sqlx.DB, cids ...string) (*MusicGroup, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT *
    FROM musicgroup
    WHERE $1 <@ cids;
    `, pq.Array(cids))

	musicgroup := &MusicGroup{}
	err := db.Get(musicgroup, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting MusicGroup (%v)", err)
	}

	members, err := musicgroup.GetMembers(db, musicgroup.ID)
	if err != nil {
		lg.Errorf("Error Scanning MembersObject (%v)", err)
	}

	image, err := musicgroup.GetImage(db, *musicgroup.ImageID)
	if err != nil {
		lg.Errorf("Error Scanning ImageObject (%v)", err)
	}

	musicgroup.Members = members
	musicgroup.Image = image

	return musicgroup, nil
}

// Create a MusicGroup
func (obj *MusicGroup) Create(db *sqlx.DB) (created *MusicGroup, err error) {
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

	var imageID IDType

	// TODO - Looks like QueryRow is being used here instead of tx.Exec
	// for the purpose of scanning ids into the returned object
	// Mixing these methods within the tx seems problematic and needs looking into
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
		return nil, fmt.Errorf("Error inserting into imageobject (%v)", err)
	}

	obj.Image.ID = imageID

	var musicgroupID IDType

	CIDs := []string{obj.CID} // pq.Array requires an array type

	err = tx.QueryRow(
		`INSERT INTO musicgroup(
            cid, cids, type, context, created_at, updated_at, name, description, email, image_id
        ) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, currval('imageobject_id_seq'))
        RETURNING id`,
		obj.CID, pq.Array(CIDs), obj.Type, obj.Context, time.Now(), time.Now(),
		obj.Name, obj.Description, obj.Email,
	).Scan(&musicgroupID)

	if err != nil {
		return nil, fmt.Errorf("Error inserting into musicgroup (%v)", err)
	}

	obj.ID = musicgroupID

	_, err = obj.UpdateMembers(db, obj.Members, tx, musicgroupID)

	if err != nil {
		return nil, fmt.Errorf("Error creating members (%v)", err)
	}

	return obj, nil
}

// Update a MusicGroup
// Note that PUT requests require id fields to be sent for sub-objects
func (obj *MusicGroup) Update(db *sqlx.DB, id IDType) (updated *MusicGroup, err error) {
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

	_, err = tx.Exec(
		`UPDATE musicgroup
         SET cid=$1, cids=array_append(cids, $2), updated_at=$3, type=$4,
         context=$5, name=$6, description=$7, email=$8
         WHERE id=$9`,
		obj.CID, obj.CID, time.Now(), obj.Type,
		obj.Context, obj.Name, obj.Description, obj.Email,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("Error updating musicgroup (%v)", err)
	}

	_, err = tx.Exec(
		`UPDATE imageobject
         SET cid=$1, type=$2, context=$3,
         updated_at=$4, content_url=$5, encoding_format=$6
         WHERE id=$7`,
		obj.Image.CID, obj.Image.Type, obj.Image.Context, time.Now(),
		obj.Image.ContentURL, obj.Image.EncodingFormat, obj.Image.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("Error updating imageobject (%v)", err)
	}

	_, err = obj.UpdateMembers(db, obj.Members, tx, id)
	if err != nil {
		return nil, fmt.Errorf("Error updating members (%v)", err)
	}

	return obj, nil
}

// Delete a MusicGroup
// TODO - Deletes remove members from the musicgroup too via CASCADE
// Potentially need to add CASCADES on other models
func (obj *MusicGroup) Delete(db *sqlx.DB, id IDType) (err error) {
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

	// Currently this also deletes all person rows from the db for members of a musicgroup
	// Eventually we might allow a musicgroup to be deleted, but for the persons to remain,
	// but for now deleting the persons simplifies testing and satisfies our use cases
	_, err = db.Exec(
		`DELETE FROM person
         WHERE id in (
                SELECT musicgroup_members.person_id
                FROM person
                INNER JOIN musicgroup_members ON musicgroup_members.person_id = person.id
                INNER JOIN musicgroup ON musicgroup_members.musicgroup_id = musicgroup.id
                WHERE musicgroup.id = $1
        )`, id,
	)
	if err != nil {
		return fmt.Errorf("Error deleting person (%v)", err)
	}

	_, err = db.Exec(
		`DELETE FROM musicgroup WHERE id=$1`, id,
	)
	if err != nil {
		return fmt.Errorf("Error deleting musicgroup (%v)", err)
	}
	return nil
}

// UpdateMembers updates the members in a group
// If member has been removed, delete person and member objects
// If person not been created yet, create the person and their member object
// If person has been created but isn't yet a member, add their member object
func (obj *MusicGroup) UpdateMembers(db *sqlx.DB, newMembers []*Person, tx *sql.Tx, id IDType) ([]*Person, error) {
	oldMembers, err := obj.GetMembers(db, id)
	if err != nil {
		return nil, fmt.Errorf("Error getting members (%v)", err)
	}

	// If member has been removed, delete person and member objects
	for _, oldMember := range oldMembers {
		var match bool
		for _, m := range newMembers {
			if m.ID == oldMember.ID {
				match = true
			}
		}

		if match == false {
			if err != nil {
				return nil, fmt.Errorf("Error deleting person (%v)", err)
			}

			_, err = tx.Exec(
				`DELETE FROM musicgroup_members
				 WHERE musicgroup_members.person_id = $1 AND musicgroup_members.musicgroup_id = $2`,
				oldMember.ID, id,
			)

			if err != nil {
				return nil, fmt.Errorf("Error deleting group (%v)", err)
			}
		}
	}

	// If person not created yet, create the person and the member object
	for _, m := range newMembers {
		var personID IDType
		if m.ID == 0 {
			personID, err = m.Create(db)

			if err != nil {
				return nil, fmt.Errorf("Error creating Person (%v)", err)
			}

			m.ID = personID
			_, err = tx.Exec(
				`INSERT INTO musicgroup_members(
					musicgroup_id, person_id, description, percentage_shares, musicgroup_admin
				) VALUES($1, $2, $3, $4, $5)`,
				id, m.ID, m.Description, m.PercentageShares, m.MusicGroupAdmin,
			)
		} else {
			// Update person
			_, err = m.Update(db, m.ID)

			if err != nil {
				return nil, fmt.Errorf("Error updating the person (%v)", err)
			}

			// If person is not yet a member, add their member object
			query := &SelectQuery{}
			query.SetQuery(`
			SELECT person.id, person.cid, person.type,
				person.context, person.created_at, person.updated_at, person.ethereum_address,
				person.given_name, person.family_name, person.email,
				musicgroup_members.description, musicgroup_members.percentage_shares, musicgroup_members.musicgroup_admin
			FROM person
			INNER JOIN musicgroup_members ON musicgroup_members.person_id = person.id
			WHERE musicgroup_members.musicgroup_id=$1 AND musicgroup_members.person_id=$2
			`, id, m.ID)

			persons := []*Person{}
			err := db.Select(&persons, query.Query(), query.Args()...)
			if err != nil {
				return nil, fmt.Errorf("Error getting MusicGroup members (%v)", err)
			}

			if len(persons) < 1 {
				_, err = tx.Exec(
					`INSERT INTO musicgroup_members(
						musicgroup_id, person_id, description, percentage_shares, musicgroup_admin
					) VALUES($1, $2, $3, $4, $5)`,
					id, m.ID, m.Description, m.PercentageShares, m.MusicGroupAdmin,
				)

				if err != nil {
					return nil, fmt.Errorf("Error adding musicgroup_member (%v)", err)
				}
			} else {
				// Update there member object if they are already a member
				_, err = tx.Exec(`
					UPDATE musicgroup_members
					SET description=$1,
					percentage_shares=$2, musicgroup_admin=$3
					WHERE musicgroup_members.musicgroup_id=$4 AND musicgroup_members.person_id=$5
					RETURNING id`,
					m.Description, m.PercentageShares, m.MusicGroupAdmin, id, m.ID,
				)

				if err != nil {
					return nil, fmt.Errorf("Error updating musicgroup_member (%v)", err)
				}
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Error adding musicgroup_member (%v)", err)
	}

	return newMembers, nil
}

// GetMembers returns all the members of a MusicGroup
func (obj *MusicGroup) GetMembers(db *sqlx.DB, musicGroupID IDType) ([]*Person, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT person.id, person.cid, person.cids, person.type,
        person.context, person.created_at, person.updated_at, person.ethereum_address,
        person.given_name, person.family_name, person.email, person.image_id, person.payment_address,
        musicgroup_members.description, musicgroup_members.percentage_shares, musicgroup_members.musicgroup_admin
    FROM person
    INNER JOIN musicgroup_members ON musicgroup_members.person_id = person.id
    WHERE musicgroup_members.musicgroup_id=$1
    `, musicGroupID)

	persons := []*Person{}
	err := db.Select(&persons, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting MusicGroup members (%v)", err)
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

	return persons, nil
}

// GetImage gets the ImageObject for a MusicGroup
func (obj *MusicGroup) GetImage(db *sqlx.DB, imageObjectID IDType) (*ImageObject, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT imageobject.id, imageobject.cid, imageobject.type, imageobject.context,
        imageobject.created_at, imageobject.updated_at, imageobject.content_url, imageobject.encoding_format
    FROM imageobject
    WHERE id=$1
    `, imageObjectID)

	imageobject := &ImageObject{}
	err := db.Get(imageobject, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting GetImage (%v)", err)
	}

	return imageobject, err
}
