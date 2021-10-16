package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/jmoiron/sqlx"
)

// Purchase model
type Purchase struct {
	ID             IDType        `json:"id"`
	CID            string        `json:"cid"`
	Type           string        `json:"@type"`
	Context        string        `json:"@context"`
	CreatedAt      time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time     `json:"updatedAt" db:"updated_at"`
	TxHash         string        `json:"txHash" db:"tx_hash"`
	Buyer          *Person       `json:"buyer"`
	MusicRelease   *MusicRelease `json:"musicRelease"`
	BuyerID        *IDType       `json:"-" db:"buyer_id"`
	MusicReleaseID *IDType       `json:"-" db:"musicrelease_id"`
}

// GetMusicRelease gets the GetMusicRelease for a MusicRelease
func (obj *Purchase) GetMusicRelease(db *sqlx.DB, musicReleaseID IDType) (*MusicRelease, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT musicrelease.id, musicrelease.cid, musicrelease.type,
        musicrelease.context, musicrelease.created_at, musicrelease.updated_at, musicrelease.active,
        musicrelease.description, musicrelease.date_published, musicrelease.catalog_number, musicrelease.music_release_format, musicrelease.price,
        musicrelease.record_label_id, musicrelease.release_of_id, image_id
    FROM musicrelease
    WHERE musicrelease.id=$1 AND musicrelease.active = true
    `, musicReleaseID)

	musicrelease := &MusicRelease{}
	err := db.Get(musicrelease, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting MusicRelease (%v)", err)
	}

	musicrelease.ReleaseOf, err = (&MusicAlbum{}).Get(db, *musicrelease.ReleaseOfID)
	if err != nil {
		lg.Errorf("Error getting MusicAlbum (%v)", err)
	}

	musicrelease.Image, err = musicrelease.GetImage(db, *musicrelease.ImageID)
	if err != nil {
		lg.Errorf("Error getting ImageObject (%v)", err)
	}
	return musicrelease, nil
}

// All gets all of the given Purchase objects
func (obj *Purchase) All(db *sqlx.DB, releaseID *IDType, ethereumAddress string, query *SelectQuery) ([]*Purchase, error) {
	// filter the purchases by a persons ethereum address and cid to check if purchased
	if ethereumAddress != "" && releaseID != nil {
		query.SetQuery(`
        SELECT purchase.id, purchase.cid, purchase.type, purchase.context,
            purchase.created_at, purchase.updated_at, tx_hash, buyer_id, musicrelease_id
        FROM purchase
        JOIN person ON purchase.buyer_id = person.id
        JOIN musicrelease ON purchase.musicrelease_id = musicrelease.id
        WHERE musicrelease.id = $1 AND musicrelease.active = true
        AND person.ethereum_address = $2
        `, releaseID, ethereumAddress)
	} else if ethereumAddress != "" && releaseID == nil { // filter the purchases by a persons ethereum address
		query.SetQuery(`
        SELECT purchase.id, purchase.cid, purchase.type, purchase.context,
            purchase.created_at, purchase.updated_at, tx_hash, buyer_id, musicrelease_id
        FROM purchase
		JOIN person ON purchase.buyer_id = person.id
		JOIN musicrelease ON purchase.musicrelease_id = musicrelease.id
        WHERE musicrelease.active = true AND person.ethereum_address = $1
        `, ethereumAddress)
	} else {
		query.SetQuery(`
        SELECT purchase.id, purchase.cid, purchase.type, purchase.context,
            purchase.created_at, purchase.updated_at, tx_hash, buyer_id, musicrelease_id
		FROM purchase
		JOIN musicrelease ON purchase.musicrelease_id = musicrelease.id
        WHERE musicrelease.active = true
        `)
	}

	purchases := []*Purchase{}
	err := db.Select(&purchases, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting all Purchases (%v)", err)
	}

	for _, purchase := range purchases {
		purchase.MusicRelease, err = purchase.GetMusicRelease(db, *purchase.MusicReleaseID)
		if err != nil {
			return nil, err
		}

		purchase.Buyer, err = (&Person{}).Get(db, *purchase.BuyerID, "")
		if err != nil {
			return nil, err
		}
	}
	return purchases, nil
}

// Get Purchase by ID
func (obj *Purchase) Get(db *sqlx.DB, id IDType) (*Purchase, error) {
	query := &SelectQuery{}
	query.SetQuery(`
        SELECT purchase.id, purchase.cid, purchase.type,
            purchase.context, purchase.created_at, purchase.updated_at, tx_hash, buyer_id, musicrelease_id
		FROM purchase
		JOIN musicrelease ON purchase.musicrelease_id = musicrelease.id
        WHERE purchase.id=$1 AND musicrelease.active = true
    `, id)

	purchase := &Purchase{}
	err := db.Get(purchase, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting Purchase (%v)", err)
	}

	purchase.MusicRelease, err = purchase.GetMusicRelease(db, *purchase.MusicReleaseID)
	if err != nil {
		return nil, err
	}

	purchase.Buyer, err = (&Person{}).Get(db, *purchase.BuyerID, "")
	if err != nil {
		return nil, err
	}
	return purchase, nil
}

// Create a Purchase
func (obj *Purchase) Create(db *sqlx.DB) (created *Purchase, err error) {
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

	var purchaseID IDType

	err = tx.QueryRow(
		`INSERT INTO purchase(
         cid, type, context, created_at, updated_at, tx_hash, buyer_id, musicrelease_id
     ) VALUES($1, $2, $3, $4, $5, $6, $7, $8)
     RETURNING id`,
		obj.CID, obj.Type, obj.Context, time.Now(), time.Now(),
		obj.TxHash, obj.Buyer.ID, obj.MusicRelease.ID,
	).Scan(&purchaseID)

	if err != nil {
		return nil, fmt.Errorf("Error creating Purchase (%v)", err)
	}

	obj.ID = purchaseID

	return obj, nil
}

// Update a Purchase
func (obj *Purchase) Update(db *sqlx.DB) error {
	return errors.New("Not implemented")
}

// Delete a Purchase
func (obj *Purchase) Delete(db *sqlx.DB, id IDType) (err error) {
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
		`DELETE FROM purchase WHERE id=$1`, id,
	)
	if err != nil {
		return fmt.Errorf("Error deleting purchase (%v)", err)
	}
	return nil
}
